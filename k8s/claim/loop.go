package claim

import (
	"context"
	"errors"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	labelKeyType             = "type"
	labelValServicePlanClaim = "service-plan-claim"
)

var (
	claimLabelSelector = labels.SelectorFromSet(labels.Set(map[string]string{
		labelKeyType: labelValServicePlanClaim,
	}))
	errLoopStopped = errors.New("loop has been stopped")
	errWatchClosed = errors.New("watch closed")
)

// StartControlLoop starts an infinite loop that receives on watcher.ResultChan() and takes action on each change in service plan claims. It's intended to be called in a goroutine. Call watcher.Stop() to stop this loop
func StartControlLoop(
	ctx context.Context,
	iface Interactor,
	cmNamespacer kcl.ConfigMapsNamespacer,
	lookup k8s.ServiceCatalogLookup,
	lifecycler *mode.Lifecycler,
) error {

	// start up the watcher so that events build up on the channel while we're listing events (which happens below)
	watcher := iface.Watch(api.ListOptions{LabelSelector: claimLabelSelector})
	ch := watcher.ResultChan()
	defer watcher.Stop()

	// iterate through all existing claims before streaming them
	claimList, err := iface.List(api.ListOptions{LabelSelector: claimLabelSelector})
	if err != nil {
		return err
	}

	for _, wrapper := range claimList.Claims {
		evt := &Event{claim: wrapper, operation: watch.Added}
		receiveEvent(ctx, evt, iface, cmNamespacer, lookup, lifecycler)
	}

	for {
		select {
		case evt, open := <-ch:
			// if the watch channel was closed, fail. this if statement is this is semantically equivalent to if evt == nil {...}
			if !open {
				return errWatchClosed
			}
			logger.Debugf("received event %s", *evt.claim.Claim)
			receiveEvent(ctx, evt, iface, cmNamespacer, lookup, lifecycler)
		case <-ctx.Done():
			logger.Debugf("loop has been stopped")
			return errLoopStopped
		}
	}
}

func receiveEvent(
	ctx context.Context,
	evt *Event,
	iface Interactor,
	cmNamespacer kcl.ConfigMapsNamespacer,
	lookup k8s.ServiceCatalogLookup,
	lifecycler *mode.Lifecycler,
) {
	nextAction, err := evt.nextAction()
	if isNoNextActionErr(err) {
		logger.Debugf("received event that has no next action (%s), skipping", err)
		return
	} else if err != nil {
		logger.Debugf("unknown error when determining the next action to make on the claim (%s)", err)
		return
	}

	claimUpdateCh := make(chan state.Update)
	wrapper := evt.claim

	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go nextAction(cancelCtx, evt, cmNamespacer, lookup, lifecycler, claimUpdateCh)

	for {
		select {
		case claimUpdate := <-claimUpdateCh:
			state.UpdateClaim(wrapper.Claim, claimUpdate)

			w, err := iface.Update(wrapper)
			if err != nil {
				logger.Errorf("error updating claim %s (%s), stopping", wrapper.Claim.ClaimID, err)
				return
			}

			// if the claim update represents a terminal state, the loop should terminate
			if claimUpdate.IsTerminal() {
				return
			}
			wrapper = w
		case <-ctx.Done():
			return
		}
	}
}
