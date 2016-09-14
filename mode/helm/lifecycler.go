package helm

import (
	"context"

	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// newLifecycler creates a new mode.Lifecycler that's backed by a Tiller instance accessible by iface
func newLifecycler(
	ctx context.Context,
	chart *chart.Chart,
	installNS string,
	provBehavior ProvisionBehavior,
	creatorDeleter ReleaseCreatorDeleter,
	cmNamespacer kcl.ConfigMapsNamespacer,
) (*mode.Lifecycler, error) {

	binder, err := newBinder(chart, cmNamespacer)
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
