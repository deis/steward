package claim

import (
	"fmt"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type errNoNextAction struct {
	evt *Event
}

func (e errNoNextAction) Error() string {
	claim := e.evt.claim.Claim
	return fmt.Sprintf(
		"no next action for operation %s with event status %s, action %s",
		e.evt.operation,
		claim.Status,
		claim.Action,
	)
}

func isNoNextActionErr(e error) bool {
	_, ok := e.(errNoNextAction)
	return ok
}

type nextFunc func(
	context.Context,
	*Event,
	kcl.ConfigMapsNamespacer,
	k8s.ServiceCatalogLookup,
	*mode.Lifecycler,
	chan<- claimUpdate,
)

// composes a bunch of nextFuncs together to make one
func compoundNextFunc(funcs ...nextFunc) nextFunc {
	return func(
		ctx context.Context,
		evt *Event,
		cmns kcl.ConfigMapsNamespacer,
		scl k8s.ServiceCatalogLookup,
		lc *mode.Lifecycler,
		ch chan<- claimUpdate) {
		for _, fn := range funcs {
			// before calling the next function, check to see if we're done
			select {
			case <-ctx.Done():
				return
			default:
			}
			fn(ctx, evt, cmns, scl, lc, ch)
		}
	}
}

// Event represents the event that a service plan claim has changed in kubernetes. It implements fmt.Stringer
type Event struct {
	claim     *ServicePlanClaimWrapper
	operation watch.EventType
}

func eventFromConfigMapEvent(raw watch.Event) (*Event, error) {
	configMap, ok := raw.Object.(*api.ConfigMap)
	if !ok {
		return nil, errNotAConfigMap
	}
	claimWrapper, err := servicePlanClaimWrapperFromConfigMap(configMap)
	if err != nil {
		return nil, err
	}
	return &Event{
		claim:     claimWrapper,
		operation: raw.Type,
	}, nil
}

func (e Event) toConfigMap() *api.ConfigMap {
	return e.claim.toConfigMap()
}

// String is the fmt.Stringer interface implementation
func (e Event) String() string {
	return fmt.Sprintf("%s %s", e.operation, *e.claim)
}

func (e *Event) nextAction() (nextFunc, error) {
	claim := e.claim.Claim

	// if event was ADDED and action is provison, next action is processProvision
	if e.operation == watch.Added && mode.StringIsAction(claim.Action, mode.ActionProvision) {
		return nextFunc(processProvision), nil
	}
	// if event was ADDED and action is bind, next action is processBind
	if e.operation == watch.Added && mode.StringIsAction(claim.Action, mode.ActionBind) {
		return nextFunc(processBind), nil
	}
	// if event was ADDED and action is create, next action is processProvision+processBind
	if e.operation == watch.Added && mode.StringIsAction(claim.Action, mode.ActionCreate) {
		return compoundNextFunc(processProvision, processBind), nil
	}
	// if event was MODIFIED, status is provisioned and action is deprovision, next action is processDeprovision
	if e.operation == watch.Modified &&
		mode.StringIsStatus(claim.Status, mode.StatusProvisioned) &&
		mode.StringIsAction(claim.Action, mode.ActionDeprovision) {
		return nextFunc(processDeprovision), nil
	}
	// if event was MODIFIED, status is provisioned and action is bind, next action is processBind
	if e.operation == watch.Modified &&
		mode.StringIsStatus(claim.Status, mode.StatusProvisioned) &&
		mode.StringIsAction(claim.Action, mode.ActionBind) {
		return nextFunc(processBind), nil
	}
	// if event was MODIFIED, status is bound and action is delete, next action is processUnbind+processDeprovision
	if e.operation == watch.Modified &&
		mode.StringIsStatus(claim.Status, mode.StatusBound) &&
		mode.StringIsAction(claim.Action, mode.ActionDelete) {
		return compoundNextFunc(processUnbind, processDeprovision), nil
	}
	// if event was MODIFIED, status is bound and action is unbind, next action is processUnbind
	if e.operation == watch.Modified &&
		mode.StringIsStatus(claim.Status, mode.StatusBound) &&
		mode.StringIsAction(claim.Action, mode.ActionUnbind) {
		return nextFunc(processUnbind), nil
	}
	// if event was MODIFIED, status is unbound and action is deprovision, next action is processDeprovision
	if e.operation == watch.Modified &&
		mode.StringIsStatus(claim.Status, mode.StatusUnbound) &&
		mode.StringIsAction(claim.Action, mode.ActionDeprovision) {
		return nextFunc(processDeprovision), nil
	}
	return nil, errNoNextAction{evt: e}
}
