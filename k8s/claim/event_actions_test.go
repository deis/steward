package claim

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/fake"
	"github.com/pborman/uuid"
)

const (
	waitDur = 100 * time.Millisecond
)

var (
	ctx = context.Background()
)

func TestNewClaimUpdate(t *testing.T) {
	// claims that should be marked stop
	stopClaims := []mode.ServicePlanClaim{
		getClaimWithStatus(mode.ActionProvision, mode.StatusFailed),
		getClaimWithStatus(mode.ActionProvision, mode.StatusProvisioned),
		getClaimWithStatus(mode.ActionProvision, mode.StatusBound),
		getClaimWithStatus(mode.ActionProvision, mode.StatusUnbound),
		getClaimWithStatus(mode.ActionProvision, mode.StatusDeprovisioned),
	}
	for i, claim := range stopClaims {
		update := state.NewUpdate(mode.Status(claim.Status), "", mode.EmptyJSONObject())
		if !update.IsTerminal() {
			t.Fatalf("update %d for claim %s was not a terminate claim", i, claim)
		}
	}

	// normal claim
	claim := getClaimWithStatus(mode.ActionProvision, mode.StatusBinding)
	update := state.NewUpdate(mode.Status(claim.Status), "", mode.EmptyJSONObject())
	assert.False(t, update.IsTerminal(), "claim was marked stop when it shouldn't have been")
}

func TestNewErrClaimUpdate(t *testing.T) {
	err := errors.New("test error")
	update := state.ErrUpdate(err, mode.EmptyJSONObject())
	assert.True(t, update.IsTerminal(), "new claim wasn't marked stop")
	assert.Equal(t, update.NewStatus, mode.StatusFailed, "resulting status")
	assert.Equal(t, err.Error(), update.NewStatusDescription, "resulting status description")
}

func getCatalogFromEvents(evts ...*Event) k8s.ServiceCatalogLookup {
	ret := k8s.NewServiceCatalogLookup(nil)
	for _, evt := range evts {
		ret.Set(&k8s.ServiceCatalogEntry{
			Info: mode.ServiceInfo{ID: evt.claim.Claim.ServiceID},
			Plan: mode.ServicePlan{ID: evt.claim.Claim.PlanID},
		})
	}
	return ret
}

func TestGetService(t *testing.T) {
	claim := getClaim(mode.ActionProvision)
	catalog := k8s.NewServiceCatalogLookup(nil)
	svc, err := getService(claim, catalog)
	assert.True(t, isNoSuchServiceAndPlanErr(err), "returned error was not a errNoSuchServiceAndPlan")
	assert.Nil(t, svc, "returned service")
	catalog = k8s.NewServiceCatalogLookup([]*k8s.ServiceCatalogEntry{
		{
			Info: mode.ServiceInfo{ID: claim.ServiceID},
			Plan: mode.ServicePlan{ID: claim.PlanID},
		},
	})
	svc, err = getService(claim, catalog)
	assert.NoErr(t, err)
	assert.NotNil(t, svc, "returned service")
	claim.ServiceID = "doesnt-exist"
	svc, err = getService(claim, catalog)
	assert.True(t, isNoSuchServiceAndPlanErr(err), "returned error was not a errNoSuchServiceAndPlan")
	assert.Nil(t, svc, "returned service")
}

func TestProcessProvisionServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionProvision))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processProvision(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not returned as terminal")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessProvisionServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionProvision))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	provisioner := &fake.Provisioner{
		Resp: &mode.ProvisionResponse{
			Extra: mode.JSONObject(map[string]string{
				uuid.New(): uuid.New(),
			}),
		},
	}
	lifecycler := &mode.Lifecycler{
		Provisioner: provisioner,
	}
	go processProvision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// provisioning status
	select {
	case claimUpdate := <-ch:
		assert.False(t, claimUpdate.IsTerminal(), "update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusProvisioning, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// provisioned status
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusProvisioned, "new status")
		assert.Equal(t, len(provisioner.Provisioned), 1, "number of provision calls")
		provCall := provisioner.Provisioned[0]
		assert.Equal(t, provCall.Req.ServiceID, evt.claim.Claim.ServiceID, "service ID")
		assert.Equal(t, provCall.Req.PlanID, evt.claim.Claim.PlanID, "plan ID")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessBindServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not returned terminal when it should have been")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessBindServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, nil, catalogLookup, nil, ch)

	// binding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, claimUpdate.IsTerminal(), "claim update was returned terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusBinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not returned terminal when it should have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusFailed, "new status")
		assert.Equal(t, claimUpdate.NewStatusDescription, errMissingInstanceID.Error(), "new status description")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessBindInstanceIDFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	evt.claim.Claim.InstanceID = uuid.New()
	catalogLookup := getCatalogFromEvents(evt)
	binder := &fake.Binder{
		Res: &mode.BindResponse{
			Creds: mode.JSONObject(map[string]string{
				"cred1": uuid.New(),
				"cred2": uuid.New(),
			}),
		},
	}
	lifecycler := &mode.Lifecycler{
		Binder: binder,
	}
	cmNamespacer := k8s.NewFakeConfigMapsNamespacer()
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, cmNamespacer, catalogLookup, lifecycler, ch)

	// binding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, claimUpdate.IsTerminal(), "claim update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusBinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// bound status
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusBound, "new status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// check the lifecycler
	assert.Equal(t, len(binder.Binds), 1, "number of bind calls")
	bindCall := binder.Binds[0]
	assert.Equal(t, bindCall.InstanceID, evt.claim.Claim.InstanceID, "instance ID")
	assert.Equal(t, bindCall.Req.ServiceID, evt.claim.Claim.ServiceID, "service ID")
	assert.Equal(t, bindCall.Req.PlanID, evt.claim.Claim.PlanID, "plan ID")

	// TODO: check the config maps namespacer

}

func TestProcessUnbindServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionUnbind))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	go processUnbind(ctx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not marked terminal when it should have been")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessUnbindServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processUnbind(cancelCtx, evt, nil, catalogLookup, nil, ch)

	// unbinding status
	select {
	case claimUpdate := <-ch:
		assert.False(t, claimUpdate.IsTerminal(), "claim update was marked terminal when it shouldn't have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusUnbinding, "new status")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.NewStatusDescription, errMissingInstanceID.Error(), "new status description")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessUnbindInstanceIDFound(t *testing.T) {
	t.Skip("TODO")
}

func TestProcessDeprovisionServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionDeprovision))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not marked terminal when it should have been")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessDeprovisionServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionDeprovision))
	catalogLookup := getCatalogFromEvents(evt)
	deprovisioner := &fake.Deprovisioner{}
	lifecycler := &mode.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	ch := make(chan state.Update)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// deprovisioning status
	select {
	case claimUpdate := <-ch:
		assert.False(t, claimUpdate.IsTerminal(), "claim update was marked terminal when it should have been")
		assert.Equal(t, claimUpdate.NewStatus, mode.StatusDeprovisioning, "new status")
		assert.Equal(t, len(deprovisioner.Deprovisions), 0, "number of deprovision calls")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.True(t, claimUpdate.IsTerminal(), "claim update was not marked terminal when it should have been")
		assert.Equal(t, claimUpdate.NewStatusDescription, errMissingInstanceID.Error(), "new status description")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestDeprovisionInstanceIDFound(t *testing.T) {
	t.Skip("TODO")
}
