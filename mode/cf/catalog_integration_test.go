// +build integration

package cf

import (
	"testing"

	"github.com/arschles/assert"
)

func TestCFCataloger(t *testing.T) {
	services, err := testCataloger.List()
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	expectedServiceCount := 3
	expectedPlanCount := 4
	assert.Equal(t, len(services), expectedServiceCount, "service count")
	for _, service := range services {
		assert.Equal(t, len(service.Plans), expectedPlanCount, "plan count")
	}
}
