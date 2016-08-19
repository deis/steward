package claim

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
)

func TestEventToConfigMap(t *testing.T) {
	evt := Event{
		claim: &ServicePlanClaimWrapper{
			Claim: &mode.ServicePlanClaim{
				TargetName: "test",
				ServiceID:  "testsvc",
				PlanID:     "testplan",
				ClaimID:    uuid.New(),
				Action:     "create",
			},
		},
	}
	configMap := evt.toConfigMap()
	assert.NoErr(t, matchClaimToMap(evt.claim.Claim, configMap.Data))
}

func TestNextAction(t *testing.T) {
	// test the state machine here
	t.Skip("TODO")
}
