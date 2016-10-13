package k8s

import (
	"testing"

	"github.com/arschles/assert"
)

func TestStringIsStatus(t *testing.T) {
	assert.True(
		t,
		StringIsStatus("teststatus", ServicePlanClaimStatus("teststatus")),
		"equivalent statuses were not reported equal",
	)
}

func TestServicePlanClaimStringer(t *testing.T) {
	assert.Equal(t, ServicePlanClaimStatus("teststatus").String(), "teststatus", "status string")
}
