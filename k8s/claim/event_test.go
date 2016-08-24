package claim

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
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
	t.Run("ADDED event, action=provision", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionProvision))
		evt.operation = watch.Added
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("ADDED event, action=bind", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionBind))
		evt.operation = watch.Added
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("ADDED event, action=create", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionCreate))
		evt.operation = watch.Added
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=deprovision, status=provisioned", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionDeprovision))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = mode.StatusProvisioned.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=bind, status=provisioned", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionBind))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = mode.StatusProvisioned.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=delete, status=bound", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionDelete))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = mode.StatusBound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=unbind, status=bound", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionUnbind))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = mode.StatusBound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})

	t.Run("MODIFIED event, action=deprovision, status=unbound", func(t *testing.T) {
		evt := getEvent(getClaim(mode.ActionDeprovision))
		evt.operation = watch.Modified
		evt.claim.Claim.Status = mode.StatusUnbound.String()
		_, err := evt.nextAction()
		assert.NoErr(t, err)
	})
}

func TestCompoundNextFunc(t *testing.T) {
	i := 0
	j := 0
	nf1 := func(
		ctx context.Context,
		evt *Event,
		cmns kcl.ConfigMapsNamespacer,
		scl k8s.ServiceCatalogLookup,
		lc mode.Lifecycler,
		ch chan<- claimUpdate) {
		select {
		case <-ctx.Done():
		default:
		}
		i++
	}
	nf2 := func(
		ctx context.Context,
		evt *Event,
		cmns kcl.ConfigMapsNamespacer,
		scl k8s.ServiceCatalogLookup,
		lc mode.Lifecycler,
		ch chan<- claimUpdate) {
		select {
		case <-ctx.Done():
		default:
		}
		j++
	}

	bgCtx := context.Background()

	// make both functions get called
	nf := compoundNextFunc(nf1, nf2)
	cancelCtx1, cancelFn1 := context.WithCancel(bgCtx)
	defer cancelFn1()
	nf(cancelCtx1, nil, nil, k8s.ServiceCatalogLookup{}, nil, nil)
	assert.Equal(t, i, 1, "number of times 1st function was called")
	assert.Equal(t, j, 1, "number of times 2nd function was called")

	// make no functions get called
	i = 0
	j = 0
	cancelCtx2, cancelFn2 := context.WithCancel(bgCtx)
	cancelFn2()
	nf(cancelCtx2, nil, nil, k8s.ServiceCatalogLookup{}, nil, nil)
	assert.Equal(t, i, 0, "number of times 1st function was called")
	assert.Equal(t, j, 0, "number of times 2nd function was called")
}
