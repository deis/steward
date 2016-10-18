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
	asyncDeprovisionRespOperationKey = "deprovision-resp-operation"
	asyncDeprovisionPollStateKey     = "deprovision-poll-state"
	asyncDeprovisionPollCountKey     = "deprovision-poll-count"
)

func pollDeprovisionState(
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
			// When deprovisioning, treat "gone" as success
			return mode.LastOperationStateSucceeded
		}

		// If maxAsyncDuration has been exceeded
		if time.Since(startTime) > maxAsyncDuration {
			select {
			case claimCh <- state.FullUpdate(
				k8s.StatusFailed,
				"asynchronous deprovisionining has exceeded the one hour allotted; service state is unknown",
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
			k8s.StatusDeprovisioningAsync,
			"polling for asynchronous deprovisionining",
			instanceID,
			"",
			mode.JSONObject(map[string]string{
				asyncDeprovisionRespOperationKey: operation,
				asyncDeprovisionPollStateKey:     pollState.String(),
				asyncDeprovisionPollCountKey:     strconv.Itoa(pollNum),
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
				// After threee consecutive polling errors, we'll consider deprovisioning failed
				select {
				case claimCh <- state.FullUpdate(
					k8s.StatusFailed,
					"polling for asynchronous depprovisionining has failed (repeatedly); service state is unknown",
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
