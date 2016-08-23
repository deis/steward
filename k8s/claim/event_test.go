package claim

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"k8s.io/kubernetes/pkg/watch"
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

func TestIsNoNextActionErr(t *testing.T) {
	e1 := errors.New("test err")
	assert.False(t, isNoNextActionErr(e1), "error was reported an errNoNextAction, but wasn't")
	e2 := errNoNextAction{evt: getEvent(getClaim(mode.ActionProvision))}
	assert.True(t, isNoNextActionErr(e2), "error was not reported an errNoNextAction, but was")
}

func TestNextAction(t *testing.T) {
	// ADDED an event to provision
	evt := getEvent(getClaim(mode.ActionProvision))
	evt.operation = watch.Added
	_, err := evt.nextAction()
	assert.NoErr(t, err)

	// ADDED an event to bind
	evt = getEvent(getClaim(mode.ActionBind))
	evt.operation = watch.Added
	_, err = evt.nextAction()
	assert.NoErr(t, err)
}
