// +build integration

package cf

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
)

const (
	// cf-sample-broker isn't picky about inputs:
	fakeInstanceID = "fake-instance-id"
	fakeServiceID  = "fake-service-id"
	fakePlanID     = "fake-plan-id"
	fakeBindingID  = "fake-binding-id"
)

func TestCfProvision(t *testing.T) {
	resp, err := testLifecycler.Provision(fakeInstanceID, &mode.ProvisionRequest{
		ServiceID: fakeServiceID,
		PlanID:    fakePlanID,
	})
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	assert.Equal(t, resp, &mode.ProvisionResponse{
		Operation: "create",
	}, "provision response")
}

func TestCfBind(t *testing.T) {
	resp, err := testLifecycler.Bind(fakeInstanceID, fakeBindingID, &mode.BindRequest{
		ServiceID: fakeServiceID,
		PlanID:    fakePlanID,
	})
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	assert.Equal(t, len(resp.Creds), 10, "credentials count")
}

func TestCfUnbind(t *testing.T) {
	err := testLifecycler.Unbind(fakeInstanceID, fakeBindingID, &mode.UnbindRequest{
		ServiceID: fakeServiceID,
		PlanID:    fakePlanID,
	})
	// Unbind returns no result except for any error that occurred...
	assert.NoErr(t, err)
}

func TestCfDeprovision(t *testing.T) {
	resp, err := testLifecycler.Deprovision(fakeInstanceID, &mode.DeprovisionRequest{
		ServiceID: fakeServiceID,
		PlanID:    fakePlanID,
	})
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	assert.Equal(t, resp, &mode.DeprovisionResponse{
		Operation: "",
	}, "deprovision response")
}
