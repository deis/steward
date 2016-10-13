package k8s

import (
	"testing"

	"github.com/arschles/assert"
)

func TestStringIsServicePlanClaimAction(t *testing.T) {
	assert.True(
		t,
		StringIsServicePlanClaimAction("testaction", ServicePlanClaimAction("testaction")),
		"'testaction' was not reported as the equivalent ServicePlanClaimAction",
	)
}

func TestServicePlanClaimActionStringer(t *testing.T) {
	assert.Equal(t, ServicePlanClaimAction("testaction").String(), "testaction", "string value")
}
