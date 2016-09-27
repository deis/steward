// +build integration

package cmd

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/test-utils/k8s"
)

func TestCmdCataloger(t *testing.T) {
	clientset, err := k8s.GetClientset()
	assert.NoErr(t, err)
	cataloger, _, err := GetComponents(clientset)
	assert.NoErr(t, err)
	services, err := cataloger.List()
	assert.NoErr(t, err)
	// Compare to known results from cmd-sample-broker...
	expectedServiceCount := 3
	expectedPlanCount := 4
	assert.Equal(t, len(services), expectedServiceCount, "service count")
	for _, service := range services {
		assert.Equal(t, len(service.Plans), expectedPlanCount, "plan count")
	}
}
