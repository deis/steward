package helm

import (
	"fmt"

	"github.com/deis/steward/mode"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	chartTmpDirPrefix = "steward-bind-chart-dl"
	tmpChartName      = "tmpchart"
)

type binder struct {
	cmInfos      []cmNamespaceAndName
	cmNamespacer kcl.ConfigMapsNamespacer
}

func dataFieldKey(cm *api.ConfigMap, key string) string {
	return fmt.Sprintf("%s-%s-%s", cm.Namespace, cm.Name, key)
}

func (b binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	resp := &mode.BindResponse{
		Creds: mode.JSONObject(map[string]string{}),
	}

	// use b.cmInfo to try and find all the listed ConfigMaps in k8s. use the data from each ConfigMap to fill in the bind response's Data field
	if err := rangeConfigMaps(b.cmNamespacer, b.cmInfos, func(cm *api.ConfigMap) error {
		for key, val := range cm.Data {
			resp.Creds[dataFieldKey(cm, key)] = val
		}
		return nil
	}); err != nil {
		logger.Errorf("ranging over helm chart credential config maps (%s)", err)
		return nil, err
	}
	return resp, nil
}

// newBinder returns a Tiller-backed mode.Binder
func newBinder(chart *chart.Chart, cmNamespacer kcl.ConfigMapsNamespacer) (mode.Binder, error) {
	cmInfos, err := getStewardConfigMapInfo(chart.Values)
	if err != nil {
		logger.Errorf("getting steward config map info (%s)", err)
		return nil, err
	}
	logger.Debugf("got config map infos for helm chart %s", cmInfos)
	// parse the values file for steward-specific config map info
	return binder{cmInfos: cmInfos, cmNamespacer: cmNamespacer}, nil
}
