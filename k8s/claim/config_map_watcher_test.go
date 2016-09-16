package claim

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/watch"
)

func configMapData() map[string]string {
	return map[string]string{
		"target-name":      "testtarget",
		"target-namespace": "testnamespace",
		"service-id":       "testsvc",
		"plan-id":          "testplan",
		"claim-id":         uuid.New(),
		"action":           "create",
	}
}

type errClaimMapMismatch struct {
	name     string
	claimVal string
	mapVal   string
}

func (e errClaimMapMismatch) Error() string {
	return fmt.Sprintf("claim %s value (%s) doesn't match map %s value (%s)", e.name, e.claimVal, e.name, e.mapVal)
}

func matchClaimToMap(claim *mode.ServicePlanClaim, data map[string]string) error {
	if claim.TargetName != data["target-name"] {
		return errClaimMapMismatch{name: "target name", claimVal: claim.TargetName, mapVal: data["target-name"]}
	}
	if claim.ServiceID != data["service-id"] {
		return errClaimMapMismatch{name: "service ID", claimVal: claim.ServiceID, mapVal: data["service-id"]}
	}
	if claim.PlanID != data["plan-id"] {
		return errClaimMapMismatch{name: "plan ID", claimVal: claim.PlanID, mapVal: data["plan-id"]}
	}
	if claim.ClaimID != data["claim-id"] {
		return errClaimMapMismatch{name: "claim ID", claimVal: claim.ClaimID, mapVal: data["claim-id"]}
	}
	if claim.Action != data["action"] {
		return errClaimMapMismatch{name: "action", claimVal: claim.Action, mapVal: data["action"]}
	}
	return nil
}

func TestConfigMapWatcher(t *testing.T) {
	ctx := context.Background()
	const waitDur = 100 * time.Millisecond
	ifaces := []*watch.FakeWatcher{
		watch.NewFake(),
		watch.NewFake(),
	}
	cancelCtx, cancelFn := context.WithCancel(ctx)
	i := 0
	watcher := newConfigMapWatcher(cancelCtx, func() (watch.Interface, error) {
		ret := ifaces[i]
		i++
		return ret, nil
	})

	defer cancelFn()
	evtCh := watcher.ResultChan()

	// add a non-config map and test for it to be ignored
	ifaces[0].Add(&api.Pod{})
	select {
	case evt, open := <-evtCh:
		if !open {
			t.Fatalf("the event channel was closed")
		}
		t.Fatalf("received an event (%s) when not expected", *evt)
	case <-time.After(waitDur):
	}

	// add a config map and expect it to be sent on the channel
	name := uuid.New()
	data := configMapData()
	ifaces[0].Add(&v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Data:       data,
	})
	select {
	case evt, open := <-evtCh:
		if !open {
			t.Fatalf("the event channel was closed")
		}
		assert.Equal(t, evt.claim.ObjectMeta.Name, name, "name")
		assert.NoErr(t, matchClaimToMap(evt.claim.Claim, data))
	case <-time.After(waitDur):
		t.Fatalf("didn't find event after %s", waitDur)
	}

	// stop the first watch interface interface and test to ensure the second is opened
	ifaces[0].Stop()

	name = uuid.New()
	data = configMapData()
	ifaces[1].Add(&v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name},
		Data:       data,
	})
	select {
	case evt := <-evtCh:
		assert.Equal(t, evt.claim.ObjectMeta.Name, name, "name")
		assert.NoErr(t, matchClaimToMap(evt.claim.Claim, data))
	case <-time.After(waitDur):
		t.Fatalf("didn't find event after %s", waitDur)
	}

}

func TestEventFromConfigMapEvent(t *testing.T) {
	rawEvt := watch.Event{
		Type:   watch.Added,
		Object: &api.Pod{},
	}
	evt, err := eventFromConfigMapEvent(rawEvt)
	assert.Nil(t, evt, "returned *Event")
	assert.Err(t, errNotAConfigMap, err)

	rawEvt = watch.Event{
		Type: watch.Added,
		Object: &v1.ConfigMap{
			Data: map[string]string{},
		},
	}
	evt, err = eventFromConfigMapEvent(rawEvt)
	assert.Nil(t, evt, "returned *Event")
	assert.True(t, err != nil, "error was nil when expected non-nil")

	data := configMapData()
	rawEvt = watch.Event{
		Type:   watch.Added,
		Object: &v1.ConfigMap{Data: data},
	}
	evt, err = eventFromConfigMapEvent(rawEvt)
	assert.NotNil(t, evt, "returned *Event")
	assert.Nil(t, err, "error")
	assert.Equal(t, evt.operation, rawEvt.Type, "event type")
	wrapper := evt.claim
	assert.NotNil(t, wrapper, "claim wrapper")
	claim := wrapper.Claim
	assert.NotNil(t, claim, "wrapped claim")
	assert.NoErr(t, matchClaimToMap(claim, data))
}

func TestWatchLoop(t *testing.T) {

}
