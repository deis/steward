package claim

import (
	"github.com/deis/steward/mode"
	"github.com/juju/loggo"
	"github.com/pborman/uuid"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

func init() {
	logger.SetLogLevel(loggo.TRACE)
}

func getEvent(claim mode.ServicePlanClaim) *Event {
	return &Event{
		claim: &ServicePlanClaimWrapper{
			Claim: &claim,
			ObjectMeta: api.ObjectMeta{
				ResourceVersion: "1",
				Name:            "testclaim",
				Namespace:       "testns",
				Labels:          map[string]string{"label-1": "label1"},
			},
		},
		operation: watch.Added,
	}
}

func getClaim(action mode.Action) mode.ServicePlanClaim {
	return mode.ServicePlanClaim{
		TargetName: "target1",
		ServiceID:  "svc1",
		PlanID:     "plan1",
		ClaimID:    uuid.New(),
		Action:     action.String(),
	}
}

func getClaimWithStatus(action mode.Action, status mode.Status) mode.ServicePlanClaim {
	cl := getClaim(action)
	cl.Status = status.String()
	return cl
}
