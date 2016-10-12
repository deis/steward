package claim

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/pborman/uuid"
	"k8s.io/client-go/1.4/pkg/watch"
)

func TestEventToConfigMap(t *testing.T) {
	evt := Event{
		claim: &k8s.ServicePlanClaimWrapper{
			Claim: &k8s.ServicePlanClaim{
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
	e2 := errNoNextAction{evt: getEvent(getClaim(k8s.ActionProvision))}
	assert.True(t, isNoNextActionErr(e2), "error was not reported an errNoNextAction, but was")
}

func TestNextAction(t *testing.T) {
	t.Run("ADDED event, action=provision", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionProvision))
		evt.operation = watch.Added
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("ADDED event, action=bind", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionBind))
		evt.operation = watch.Added
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("ADDED event, action=create", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionCreate))
		evt.operation = watch.Added
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=delete, status=provisioned", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionCreate))
		evt.claim.Claim.Status = k8s.StatusProvisioned.String()
		evt.operation = watch.Modified
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=deprovision, status=provisioned", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionDeprovision))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = k8s.StatusProvisioned.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=bind, status=provisioned", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionBind))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = k8s.StatusProvisioned.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=delete, status=bound", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionDelete))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = k8s.StatusBound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=delete, status=unbound", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionDelete))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = k8s.StatusUnbound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=unbind, status=bound", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionUnbind))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = k8s.StatusBound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=deprovision, status=unbound", func(t *testing.T) {
		evt := getEvent(getClaim(k8s.ActionDeprovision))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = k8s.StatusUnbound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})
}
