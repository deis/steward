package claim

import (
	"context"
	"fmt"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
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
	action := mode.Action(claim.Action)
	status := mode.Status(claim.Status)

	stateNoStatus := state.NewCurrentNoStatus(action, e.operation)
	stateWithStatus := state.NewCurrent(status, action, e.operation)

	nextFuncNoStatus, okNoStatus := transitionTable[stateNoStatus]
	nextFuncWithStatus, okWithStatus := transitionTable[stateWithStatus]
	if !okNoStatus && !okWithStatus {
		return nil, errNoNextAction{evt: e}
	} else if !okNoStatus {
		return nextFuncWithStatus, nil
	} else {
		return nextFuncNoStatus, nil
	}
}
