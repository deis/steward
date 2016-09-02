package helm

import (
	"context"

	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// NewLifecycler creates a new mode.Lifecycler that's backed by a Tiller instance accessible by iface
func NewLifecycler(
	ctx context.Context,
	chart *chart.Chart,
	installNS string,
	provBehavior ProvisionBehavior,
	creatorDeleter ReleaseCreatorDeleter,
	cmNamespacer kcl.ConfigMapsNamespacer,
) (*mode.Lifecycler, error) {

	binder, err := NewBinder(chart, cmNamespacer)
	if err != nil {
		return nil, err
	}
	unbinder, err := NewUnbinder(chart, cmNamespacer)
	return &mode.Lifecycler{
		Provisioner:   NewProvisioner(chart, installNS, provBehavior, creatorDeleter),
		Binder:        binder,
		Unbinder:      unbinder,
		Deprovisioner: NewDeprovisioner(chart, provBehavior, creatorDeleter),
	}, nil
}
