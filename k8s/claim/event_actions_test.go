package claim

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/fake"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	waitDur = 100 * time.Millisecond
)

var (
	ctx = context.Background()
)

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

func getEvent(claim mode.ServicePlanClaim) *Event {
	return &Event{
		claim: &ServicePlanClaimWrapper{
			Claim: &claim,
			ObjectMeta: api.ObjectMeta{
				ResourceVersion: "1",
				Name:            "testclaim",
				Namespace:       "testns",
				Labels:          map[string]string{"label-1": "label1"},
			},
		},
		operation: watch.Added,
	}
}

func getClaim(action mode.Action) mode.ServicePlanClaim {
	return mode.ServicePlanClaim{
		TargetName: "target1",
		ServiceID:  "svc1",
		PlanID:     "plan1",
		ClaimID:    uuid.New(),
		Action:     action.String(),
	}
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
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processProvision(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "returned error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessProvisionServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionProvision))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	lifecycler := fake.Lifecycler{
		Provisioner: &fake.Provisioner{},
	}
	go processProvision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// provisioning status
	select {
	case claimUpdate := <-ch:
		assert.NoErr(t, claimUpdate.err)
		assert.False(t, claimUpdate.stop, "stop boolean in claim update")
		assert.NotNil(t, claimUpdate.newClaim, "new claim")
		assert.Equal(t, claimUpdate.newClaim.Status, mode.StatusProvisioning.String(), "status")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// provisioned status
	select {
	case claimUpdate := <-ch:
		assert.NoErr(t, claimUpdate.err)
		assert.True(t, claimUpdate.stop, "stop boolean in claim update")
		assert.NotNil(t, claimUpdate.newClaim, "new claim")
		assert.Equal(t, claimUpdate.newClaim.Status, mode.StatusProvisioned.String(), "status")
		assert.Equal(t, len(lifecycler.Provisioner.Provisioned), 1, "number of provision calls")
		provCall := lifecycler.Provisioner.Provisioned[0]
		assert.Equal(t, provCall.Req.ServiceID, evt.claim.Claim.ServiceID, "service ID")
		assert.Equal(t, provCall.Req.PlanID, evt.claim.Claim.PlanID, "plan ID")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessBindServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "return error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessBindServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processBind(cancelCtx, evt, nil, catalogLookup, nil, ch)

	// binding status
	select {
	case claimUpdate := <-ch:
		assert.NoErr(t, claimUpdate.err)
		assert.False(t, claimUpdate.stop, "stop boolean in claim update")
		assert.NotNil(t, claimUpdate.newClaim, "new claim")
		assert.Equal(t, claimUpdate.newClaim.Status, mode.StatusBinding.String(), "status")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.Err(t, claimUpdate.err, errMissingInstanceID)
		assert.True(t, claimUpdate.stop, "stop boolean")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestProcessBindInstanceIDFound(t *testing.T) {
	t.Skip("TODO")
}

func TestProcessUnbindServiceNotFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionUnbind))
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan claimUpdate)
	go processUnbind(ctx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "return error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessUnbindServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionBind))
	catalogLookup := getCatalogFromEvents(evt)
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processUnbind(cancelCtx, evt, nil, catalogLookup, nil, ch)

	// unbinding status
	select {
	case claimUpdate := <-ch:
		assert.NoErr(t, claimUpdate.err)
		assert.False(t, claimUpdate.stop, "stop boolean in claim update")
		assert.NotNil(t, claimUpdate.newClaim, "new claim")
		assert.Equal(t, claimUpdate.newClaim.Status, mode.StatusUnbinding.String(), "status")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.Err(t, claimUpdate.err, errMissingInstanceID)
		assert.True(t, claimUpdate.stop, "stop boolean")
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
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "return error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessDeprovisionServiceFound(t *testing.T) {
	evt := getEvent(getClaim(mode.ActionDeprovision))
	catalogLookup := getCatalogFromEvents(evt)
	lifecycler := fake.Lifecycler{
		Deprovisioner: &fake.Deprovisioner{},
	}
	ch := make(chan claimUpdate)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go processDeprovision(cancelCtx, evt, nil, catalogLookup, lifecycler, ch)

	// deprovisioning status
	select {
	case claimUpdate := <-ch:
		assert.NoErr(t, claimUpdate.err)
		assert.False(t, claimUpdate.stop, "stop boolean in claim update")
		assert.NotNil(t, claimUpdate.newClaim, "new claim")
		assert.Equal(t, claimUpdate.newClaim.Status, mode.StatusDeprovisioning.String(), "status")
		assert.Equal(t, len(lifecycler.Deprovisioner.Deprovisions), 0, "number of deprovision calls")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}

	// missing instance ID status
	select {
	case claimUpdate := <-ch:
		assert.Err(t, claimUpdate.err, errMissingInstanceID)
		assert.True(t, claimUpdate.stop, "stop boolean")
	case <-time.After(waitDur):
		t.Fatalf("didn't receive a claim update within %s", waitDur)
	}
}

func TestDeprovisionInstanceIDFound(t *testing.T) {
	t.Skip("TODO")
}
