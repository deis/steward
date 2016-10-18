package claim

import (
	"context"
	"strconv"
	"time"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
)

const (
	asyncProvisionRespOperationKey = "provision-resp-operation"
	asyncProvisionPollStateKey     = "provision-poll-state"
	asyncProvisionPollCountKey     = "provision-poll-count"
)

func pollProvisionState(
	ctx context.Context,
	serviceID,
	planID,
	operation,
	instanceID string,
	lastOpGetter mode.LastOperationGetter,
	claimCh chan<- state.Update,
) mode.LastOperationState {

	pollNum := 0
	pollState := mode.LastOperationStateInProgress
	pollErrCount := 0

	cfg, err := getConfig()
	if err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return mode.LastOperationStateFailed
	}
	maxAsyncDuration := cfg.getMaxAsyncDuration()

	startTime := time.Now()
	for {
		if pollState == mode.LastOperationStateSucceeded || pollState == mode.LastOperationStateFailed {
			// if the polling went into success or failed state, just return that
			return pollState
		}
		if pollState == mode.LastOperationStateGone {
			// When provisioning, treat "gone" as a failure
			return mode.LastOperationStateFailed
		}

		// If maxAsyncDuration has been exceeded
		if time.Since(startTime) > maxAsyncDuration {
			select {
			case claimCh <- state.FullUpdate(
				k8s.StatusFailed,
				"asynchronous provisionining has exceeded the one hour allotted; service state is unknown",
				instanceID,
				"",
				mode.EmptyJSONObject(),
			):
			case <-ctx.Done():
			}
			return mode.LastOperationStateFailed
		}

		// otherwise continue provisioning state
		update := state.FullUpdate(
			k8s.StatusProvisioningAsync,
			"polling for asynchronous provisionining",
			instanceID,
			"",
			mode.JSONObject(map[string]interface{}{
				asyncProvisionRespOperationKey: operation,
				asyncProvisionPollStateKey:     pollState.String(),
				asyncProvisionPollCountKey:     strconv.Itoa(pollNum),
			}))
		select {
		case claimCh <- update:
		case <-ctx.Done():
		}
		resp, err := lastOpGetter.GetLastOperation(serviceID, planID, operation, instanceID)
		if err != nil {
			if pollErrCount < 3 {
				pollErrCount++
			} else {
				// After threee consecutive polling errors, we'll consider provisioning failed
				select {
				case claimCh <- state.FullUpdate(
					k8s.StatusFailed,
					"polling for asynchronous provisionining has failed (repeatedly); service state is unknown",
					instanceID,
					"",
					mode.EmptyJSONObject(),
				):
				case <-ctx.Done():
				}
				return mode.LastOperationStateFailed
			}
		} else {
			// Reset error count to zero
			pollErrCount = 0
		}
		pollNum++
		newState, err := resp.GetState()
		if err != nil {
			select {
			case claimCh <- state.ErrUpdate(err):
			case <-ctx.Done():
			}
			return mode.LastOperationStateFailed
		}
		pollState = newState
		time.Sleep(30 * time.Second)
	}
}
