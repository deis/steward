package claim

import (
	"context"
	"strconv"
	"time"

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
	for {
		if pollState == mode.LastOperationStateSucceeded || pollState == mode.LastOperationStateFailed {
			// if the polling went into success or failed state, just return that
			return pollState
		}

		// otherwise continue provisioning state
		update := state.FullUpdate(
			mode.StatusDeprovisioningAsync,
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
			select {
			case claimCh <- state.ErrUpdate(err):
			case <-ctx.Done():
			}
			return mode.LastOperationStateFailed
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
