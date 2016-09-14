package helm

import (
	"testing"

	"github.com/arschles/assert"
)

func TestCatalogerList(t *testing.T) {
	cfg := &config{
		ServiceID:          "svcID1",
		ServiceName:        "svc1",
		ServiceDescription: "this is service 1",
		PlanID:             "planID1",
		PlanName:           "plan1",
		PlanDescription:    "this is plan 1",
	}
	cat := newCataloger(cfg)
	svcs, err := cat.List()
	assert.NoErr(t, err)
	assert.Equal(t, len(svcs), 1, "number of listed services")
	svc := svcs[0]
	assert.Equal(t, svc.ID, cfg.ServiceID, "service ID")
	assert.Equal(t, svc.Name, cfg.ServiceName, "service name")
	assert.Equal(t, svc.Description, cfg.ServiceDescription, "service description")
	assert.Equal(t, len(svc.Plans), 1, "number of plans")
	plan := svc.Plans[0]
	assert.Equal(t, plan.ID, cfg.PlanID, "plan ID")
	assert.Equal(t, plan.Name, cfg.PlanName, "plan name")
	assert.Equal(t, plan.Description, cfg.PlanDescription, "plan description")
}
