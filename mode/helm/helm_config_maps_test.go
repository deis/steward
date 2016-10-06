package helm

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"k8s.io/client-go/1.4/pkg/api/v1"
)

func TestRangeConfigMaps(t *testing.T) {
	cmNames := []string{"name1", "name2"}
	cmIface := k8s.NewFakeConfigMapsInterface()
	gathered := []*v1.ConfigMap{}
	rangeConfigMaps(cmIface, cmNames, func(cm *v1.ConfigMap) error {
		gathered = append(gathered, cm)
		return nil
	})
	assert.Equal(t, len(gathered), len(cmNames), "number of gathered config maps")
}
