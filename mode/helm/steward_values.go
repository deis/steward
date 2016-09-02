package helm

import (
	"fmt"

	"github.com/ghodss/yaml"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type cmNamespaceAndName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (c cmNamespaceAndName) String() string {
	return fmt.Sprintf("namespace=%s,name=%s", c.Namespace, c.Name)
}

type stewardValues struct {
	ConfigMaps []cmNamespaceAndName `json:"stewardConfigMaps"`
}

func getStewardConfigMapInfo(chart *chart.Config) ([]cmNamespaceAndName, error) {
	var ret stewardValues
	if err := yaml.Unmarshal([]byte(chart.Raw), &ret); err != nil {
		return nil, err
	}
	return ret.ConfigMaps, nil
}
