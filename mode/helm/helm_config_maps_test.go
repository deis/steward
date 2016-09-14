package helm

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"k8s.io/kubernetes/pkg/api"
)

func TestRangeConfigMaps(t *testing.T) {
	infos := []cmNamespaceAndName{
		{Name: "name1", Namespace: "ns1"},
		{Name: "name2", Namespace: "ns2"},
	}
	namespacer := k8s.NewFakeConfigMapsNamespacer()
	gathered := []*api.ConfigMap{}
	rangeConfigMaps(namespacer, infos, func(cm *api.ConfigMap) error {
		gathered = append(gathered, cm)
		return nil
	})
	assert.Equal(t, len(gathered), len(infos), "number of gathered config maps")
}
