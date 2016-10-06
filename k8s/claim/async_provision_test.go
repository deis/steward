package claim

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/fake"
	"github.com/pborman/uuid"
)

func TestPollProvisionState(t *testing.T) {
	const (
		timeout = 1 * time.Second
	)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	serviceID := uuid.New()
	planID := uuid.New()
	operation := "testOperation"
	instanceID := uuid.New()

	var curStateMut sync.RWMutex
	curState := mode.LastOperationStateInProgress

	lastOpGetter := &fake.LastOperationGetter{
		Ret: func() *mode.GetLastOperationResponse {
			curStateMut.RLock()
			defer curStateMut.RUnlock()
			return &mode.GetLastOperationResponse{State: curState.String()}
		},
	}

	claimCh := make(chan state.Update)
	go func() {
		finalState := pollProvisionState(ctx, serviceID, planID, operation, instanceID, lastOpGetter, claimCh)
		assert.Equal(t, finalState, mode.LastOperationStateSucceeded, "final state")
	}()

	/////
	// expect a provisioning-async first. after we receive it, the last op getter will get another provisioning-async and wait to send it. we then change the current state, receive the second provisioning-async and then expect the channel to not receive anymore. the final success state is received in the return value of pollProvisionState, and it's checked in the above goroutine
	/////

	assert.NoErr(t, acceptStatus(claimCh, mode.StatusProvisioningAsync))

	curStateMut.Lock()
	curState = mode.LastOperationStateSucceeded
	curStateMut.Unlock()

	assert.NoErr(t, acceptStatus(claimCh, mode.StatusProvisioningAsync))

	select {
	case update := <-claimCh:
		t.Fatalf("received %s on claim channel, expected nothing", update)
	case <-time.After(timeout):
	}
}

func acceptStatus(claimCh <-chan state.Update, expected mode.Status) error {
	const timeout = 1 * time.Second
	select {
	case update := <-claimCh:
		if update.Status() != expected {
			return fmt.Errorf("expected status %s, got %s", expected, update.Status())
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("no status update after %s", timeout)
	}
}
