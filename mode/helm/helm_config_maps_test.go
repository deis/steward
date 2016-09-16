package helm

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"k8s.io/client-go/1.4/pkg/api/v1"
)

func TestRangeConfigMaps(t *testing.T) {
	infos := []cmNamespaceAndName{
		{Name: "name1", Namespace: "ns1"},
		{Name: "name2", Namespace: "ns2"},
	}
	namespacer := k8s.NewFakeConfigMapsNamespacer()
	gathered := []*v1.ConfigMap{}
	rangeConfigMaps(namespacer, infos, func(cm *v1.ConfigMap) error {
		gathered = append(gathered, cm)
		return nil
	})
	assert.Equal(t, len(gathered), len(infos), "number of gathered config maps")
}
