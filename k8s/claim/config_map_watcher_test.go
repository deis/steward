package claim

import (
	"fmt"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
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

func TestConfigMapWatcherStop(t *testing.T) {
	const waitDur = 100 * time.Millisecond
	iface := watch.NewFake()
	watcher := newConfigMapWatcher(iface)
	watcher.Stop()
	ch := watcher.ResultChan()
	select {
	case evt := <-ch:
		if evt != nil {
			t.Fatalf("got event %s within %s, but shouldn't have gotten anything", *evt, waitDur)
		}
	case <-time.After(waitDur):
	}
}

func TestConfigMapWatcher(t *testing.T) {
	const waitDur = 100 * time.Millisecond
	iface := watch.NewFake()
	watcher := newConfigMapWatcher(iface)
	defer watcher.Stop()
	evtCh := watcher.ResultChan()

	iface.Add(&api.Pod{})
	select {
	case evt := <-evtCh:
		t.Fatalf("received an event (%s) when not expected", *evt)
	case <-time.After(waitDur):
	}

	data := configMapData()
	iface.Add(&api.ConfigMap{Data: data})
	select {
	case evt := <-evtCh:
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
		Object: &api.ConfigMap{
			Data: map[string]string{},
		},
	}
	evt, err = eventFromConfigMapEvent(rawEvt)
	assert.Nil(t, evt, "returned *Event")
	assert.True(t, err != nil, "error was nil when expected non-nil")

	data := configMapData()
	rawEvt = watch.Event{
		Type:   watch.Added,
		Object: &api.ConfigMap{Data: data},
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
