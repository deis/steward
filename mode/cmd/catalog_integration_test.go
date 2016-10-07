// +build integration

package cmd

import (
	"testing"

	"github.com/arschles/assert"
)

func TestCmdCataloger(t *testing.T) {
	services, err := testCataloger.List()
	assert.NoErr(t, err)
	// Compare to known results from cmd-sample-broker...
	expectedServiceCount := 3
	expectedPlanCount := 4
	assert.Equal(t, len(services), expectedServiceCount, "service count")
	for _, service := range services {
		assert.Equal(t, len(service.Plans), expectedPlanCount, "plan count")
	}
}
