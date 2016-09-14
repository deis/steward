package helm

import (
	"testing"

	"github.com/arschles/assert"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func TestProvisionNoop(t *testing.T) {
	prov := newProvisioner(nil, "targetNS", ProvisionBehaviorNoop, nil)
	resp, err := prov.Provision("testInstanceID", nil)
	assert.NoErr(t, err)
	assert.Equal(t, resp.Operation, provisionedNoopOperation, "operation")
}

func TestProvisionActive(t *testing.T) {
	const (
		targetNS    = "testNamespace"
		instID      = "testInstanceID"
		releaseName = "testRelease"
	)
	chart := &chart.Chart{}
	creator := &fakeCreatorDeleter{
		createResp: &services.InstallReleaseResponse{
			Release: &release.Release{
				Name: releaseName,
			},
		},
	}
	prov := newProvisioner(chart, targetNS, ProvisionBehaviorActive, creator)
	resp, err := prov.Provision(instID, nil)
	assert.NoErr(t, err)
	assert.NotNil(t, resp, "provision response")
	assert.Equal(t, len(creator.deleteCalls), 0, "number of calls to delete")
	assert.Equal(t, len(creator.createCalls), 1, "number of calls to create")
	createCall := creator.createCalls[0]
	assert.Equal(t, createCall.installNS, targetNS, "target namespace of create call")
	assert.Equal(t, createCall.chart, chart, "installed chart")
	assert.Equal(t, resp.Extra[releaseNameKey], releaseName, "release name")
	assert.Equal(t, resp.Operation, provisionedActiveOperation, "provision response operation")
}
