package helm

import (
	"fmt"

	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

const (
	deprovisionedNoopOperation   = "deprovisioned-noop"
	deprovisionedActiveOperation = "deprovisioned-active"
)

type errMissingReleaseName struct {
	params mode.JSONObject
}

func (e errMissingReleaseName) Error() string {
	return fmt.Sprintf("missing release name in parameters %s", e.params)
}

type deprovisioner struct {
	chart        *chart.Chart
	provBehavior ProvisionBehavior
	deleter      ReleaseDeleter
}

func (d deprovisioner) Deprovision(instanceID string, dreq *mode.DeprovisionRequest) (*mode.DeprovisionResponse, error) {
	if d.provBehavior == ProvisionBehaviorNoop {
		return &mode.DeprovisionResponse{
			Operation: deprovisionedNoopOperation,
		}, nil
	}
	_, ok := dreq.Parameters[releaseNameKey]
	if !ok {
		logger.Errorf("finding the release name key")
		return nil, errMissingReleaseName{params: dreq.Parameters}
	}
	releaseName, err := dreq.Parameters.String(releaseNameKey)
	if err != nil {
		return nil, err
	}
	if _, err := d.deleter.Delete(releaseName); err != nil {
		logger.Errorf("deleting the helm chart (%s)", err)
		return nil, err
	}

	return &mode.DeprovisionResponse{
		Operation: deprovisionedActiveOperation,
	}, nil
}

// newDeprovisioner returns a new Tiller-backed mode.Deprovisioner
func newDeprovisioner(chart *chart.Chart, provBehavior ProvisionBehavior, deleter ReleaseDeleter) mode.Deprovisioner {
	return deprovisioner{
		chart:        chart,
		provBehavior: provBehavior,
		deleter:      deleter,
	}
}
