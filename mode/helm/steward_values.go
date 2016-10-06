package helm

import (
	"github.com/ghodss/yaml"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type stewardValues struct {
	ConfigMaps []string `json:"stewardConfigMaps"`
}

func getStewardConfigMapInfo(chart *chart.Config) ([]string, error) {
	var ret stewardValues
	if err := yaml.Unmarshal([]byte(chart.Raw), &ret); err != nil {
		return nil, err
	}
	return ret.ConfigMaps, nil
}
