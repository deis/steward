package helm

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func TestDeprovisionNoop(t *testing.T) {
	deprov := newDeprovisioner(nil, ProvisionBehaviorNoop, nil)
	resp, err := deprov.Deprovision("testInstanceID", nil)
	assert.NoErr(t, err)
	assert.Equal(t, resp.Operation, deprovisionedNoopOperation, "operation")
}
func TestDeprovisionActive(t *testing.T) {
	const (
		targetNS    = "testNamespace"
		instID      = "testInstanceID"
		releaseName = "testRelease"
	)
	chart := &chart.Chart{}
	deleter := &fakeCreatorDeleter{}
	deprov := newDeprovisioner(chart, ProvisionBehaviorActive, deleter)
	deprovReq := &mode.DeprovisionRequest{
		AcceptsIncomplete: true,
		Parameters:        mode.JSONObject(map[string]interface{}{releaseNameKey: releaseName}),
	}
	resp, err := deprov.Deprovision(instID, deprovReq)
	assert.NoErr(t, err)
	assert.NotNil(t, resp, "deprovision response")
	assert.Equal(t, len(deleter.createCalls), 0, "number of calls to create")
	assert.Equal(t, len(deleter.deleteCalls), 1, "number of calls to delete")
	deleteCall := deleter.deleteCalls[0]
	assert.Equal(t, deleteCall, releaseName, "target namespace of create call")
	assert.Equal(t, resp.Operation, deprovisionedActiveOperation, "release name")

}
