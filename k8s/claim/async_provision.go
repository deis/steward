package claim

import (
	"context"
	"strconv"

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
	for {
		if pollState == mode.LastOperationStateSucceeded || pollState == mode.LastOperationStateFailed {
			// if the polling went into success or failed state, just return that
			return pollState
		}

		// otherwise continue provisioning state
		update := state.FullUpdate(
			mode.StatusProvisioningAsync,
			"polling for asynchronous provisionining",
			instanceID,
			"",
			mode.JSONObject(map[string]string{
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
	}
}
