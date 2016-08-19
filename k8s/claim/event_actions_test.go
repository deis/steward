package claim

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	waitDur = 100 * time.Millisecond
)

var (
	ctx = context.Background()
	evt = &Event{
		claim: &ServicePlanClaimWrapper{
			Claim:           getClaim(),
			ResourceVersion: "1",
			OriginalName:    "testclaim",
			Labels:          map[string]string{"label-1": "label1"},
		},
		operation: watch.Added,
	}
)

func getClaim() *mode.ServicePlanClaim {
	return &mode.ServicePlanClaim{
		TargetName:      "target1",
		TargetNamespace: "targetns1",
		ServiceID:       "svc1",
		PlanID:          "plan1",
		ClaimID:         uuid.New(),
		Action:          mode.ActionProvision.String(),
	}
}

func TestGetService(t *testing.T) {
	claim := getClaim()
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

func TestProcessProvision(t *testing.T) {
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan claimUpdate)
	go processProvision(ctx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "returned error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessBind(t *testing.T) {
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan claimUpdate)
	go processBind(ctx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "return error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}

func TestProcessUnbind(t *testing.T) {
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

func TestProcessDeprovision(t *testing.T) {
	catalogLookup := k8s.NewServiceCatalogLookup(nil)
	ch := make(chan claimUpdate)
	go processDeprovision(ctx, evt, nil, catalogLookup, nil, ch)
	select {
	case claimUpdate := <-ch:
		assert.NotNil(t, claimUpdate.err, "returned claim update error")
		assert.True(t, isNoSuchServiceAndPlanErr(claimUpdate.err), "return error should have been an errNoSuchServiceAndPlan")
	case <-time.After(waitDur):
		t.Fatalf("no claim update given after %s", waitDur)
	}
}
