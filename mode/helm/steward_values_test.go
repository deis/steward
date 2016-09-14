package helm

import (
	"testing"

	"github.com/arschles/assert"
	"k8s.io/helm/pkg/chartutil"
)

func TestGetStewardConfigMapInfo(t *testing.T) {
	chart, err := chartutil.Load(alpineChartLoc())
	assert.NoErr(t, err)
	nsAndNames, err := getStewardConfigMapInfo(chart.Values)
	assert.NoErr(t, err)
	assert.Equal(t, len(nsAndNames), 1, "number of namespace-and-name pairs")
	nsAndName := nsAndNames[0]
	// (name, namespace) values in the alpine chart are expected to be ("my-creds", "default") (respectively)
	assert.Equal(t, nsAndName.Name, "my-creds", "name of config map")
	assert.Equal(t, nsAndName.Namespace, "default", "namespace of config map")
}
