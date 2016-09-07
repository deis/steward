package helm

import (
	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

type unbinder struct {
	cmInfos      []cmNamespaceAndName
	cmNamespacer kcl.ConfigMapsNamespacer
}

func (u unbinder) Unbind(instanceID, bindingID string, unbindRequest *mode.UnbindRequest) error {
	return nil
}

// newUnbinder returns a Tiller-backed mode.Unbinder
func newUnbinder(chart *chart.Chart, cmNamespacer kcl.ConfigMapsNamespacer) (mode.Unbinder, error) {
	cmInfos, err := getStewardConfigMapInfo(chart.Values)
	if err != nil {
		return nil, err
	}
	return unbinder{
		cmInfos:      cmInfos,
		cmNamespacer: cmNamespacer,
	}, nil
}
