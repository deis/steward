package claim

import (
	"context"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/fake"
	"github.com/pborman/uuid"
)

func TestPollProvisionState(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	operation := "testOperation"
	instanceID := uuid.New()
	curState := mode.LastOperationStateInProgress

	switcherCh := make(chan mode.LastOperationState)
	lastOpGetter := &fake.LastOperationGetter{
		Ret: func() *mode.GetLastOperationResponse {
			select {
			case newState := <-switcherCh:
				curState = newState
			default:
			}
			return &mode.GetLastOperationResponse{State: curState.String()}
		},
	}

	claimCh := make(chan state.Update)
	go func() {
		finalState := pollProvisionState(ctx, serviceID, planID, operation, instanceID, lastOpGetter, claimCh)
		assert.Equal(t, finalState, mode.LastOperationStateSucceeded, "final state")
	}()

	select {
	case update := <-claimCh:
		assert.Equal(t, update.Status(), mode.StatusProvisioningAsync, "claim update status")
	case <-time.After(1 * time.Second):
		t.Fatalf("got no initial update")
	}

	// TODO: deadlock here
	select {
	case switcherCh <- mode.LastOperationStateSucceeded:
	case <-time.After(1 * time.Second):
		t.Fatalf("unable to switch the last operation state")
	}

	select {
	case update := <-claimCh:
		assert.Equal(t, update.Status(), mode.StatusProvisioned, "claim update status")
	case <-time.After(1 * time.Second):
		t.Fatalf("got no second update")
	}
}
