package helm

import (
	"fmt"

	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

const (
	chartTmpDirPrefix = "steward-bind-chart-dl"
	tmpChartName      = "tmpchart"
)

type binder struct {
	cmNames []string
	cmIface v1.ConfigMapInterface
}

func dataFieldKey(cm *v1types.ConfigMap, key string) string {
	return fmt.Sprintf("%s-%s-%s", cm.Namespace, cm.Name, key)
}

func (b binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	resp := &mode.BindResponse{
		Creds: mode.JSONObject(map[string]interface{}{}),
	}

	// use b.cmInfo to try and find all the listed ConfigMaps in k8s. use the data from each ConfigMap to fill in the bind response's Data field
	if err := rangeConfigMaps(b.cmIface, b.cmNames, func(cm *v1types.ConfigMap) error {
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
func newBinder(chart *chart.Chart, cmIface v1.ConfigMapInterface) (mode.Binder, error) {
	// parse the values file for steward-specific config map info
	cmNames, err := getStewardConfigMapInfo(chart.Values)
	if err != nil {
		logger.Errorf("getting steward config map info (%s)", err)
		return nil, err
	}
	logger.Debugf("got config map names for helm chart %s", cmNames)
	return binder{
		cmNames: cmNames,
		cmIface: cmIface,
	}, nil
}
