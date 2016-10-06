package helm

import (
	"testing"

	"github.com/arschles/assert"
	"k8s.io/helm/pkg/chartutil"
)

func TestGetStewardConfigMapInfo(t *testing.T) {
	chart, err := chartutil.Load(alpineChartLoc())
	assert.NoErr(t, err)
	cmNames, err := getStewardConfigMapInfo(chart.Values)
	assert.NoErr(t, err)
	assert.Equal(t, len(cmNames), 1, "number of namespace-and-name pairs")
	cmName := cmNames[0]
	// (name, namespace) values in the alpine chart are expected to be ("my-creds", "default") (respectively)
	assert.Equal(t, cmName, "my-creds", "name of config map")
}
