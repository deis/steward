package helm

import (
	"context"

	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// newLifecycler creates a new mode.Lifecycler that's backed by a Tiller instance accessible by iface
func newLifecycler(
	ctx context.Context,
	chart *chart.Chart,
	installNS string,
	provBehavior ProvisionBehavior,
	creatorDeleter ReleaseCreatorDeleter,
	cmNamespacer v1.ConfigMapsGetter,
) (*mode.Lifecycler, error) {

	cmIface := cmNamespacer.ConfigMaps(installNS)
	binder, err := newBinder(chart, cmIface)
	if err != nil {
		return nil, err
	}
	return &mode.Lifecycler{
		Provisioner:   newProvisioner(chart, installNS, provBehavior, creatorDeleter),
		Binder:        binder,
		Unbinder:      newUnbinder(),
		Deprovisioner: newDeprovisioner(chart, provBehavior, creatorDeleter),
	}, nil
}
