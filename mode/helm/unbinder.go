package helm

import (
	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

type unbinder struct {
	cmInfos      []cmNamespaceAndName
	cmNamespacer kcl.ConfigMapsNamespacer
}

func (u unbinder) Unbind(instanceID, bindingID string, unbindRequest *mode.UnbindRequest) error {
	// iterate and delete all ConfigMaps represented in cmInfos
	if err := rangeConfigMaps(u.cmNamespacer, u.cmInfos, func(cm *api.ConfigMap) error {
		return u.cmNamespacer.ConfigMaps(cm.Namespace).Delete(cm.Name)
	}); err != nil {
		logger.Errorf("ranging over helm chart credential config maps (%s)", err)
		return err
	}
	return nil
}

// NewUnbinder returns a Tiller-backed mode.Unbinder
func NewUnbinder(chart *chart.Chart, cmNamespacer kcl.ConfigMapsNamespacer) (mode.Unbinder, error) {
	cmInfos, err := getStewardConfigMapInfo(chart.Values)
	if err != nil {
		return nil, err
	}
	return unbinder{
		cmInfos:      cmInfos,
		cmNamespacer: cmNamespacer,
	}, nil
}
