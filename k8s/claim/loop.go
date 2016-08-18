package claim

import (
	"errors"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"golang.org/x/net/context"
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
)

// StartControlLoop starts an infinite loop that receives on watcher.ResultChan() and takes action on each change in service plan claims. It's intended to be called in a goroutine. Call watcher.Stop() to stop this loop
func StartControlLoop(
	ctx context.Context,
	iface Interactor,
	cmNamespacer kcl.ConfigMapsNamespacer,
	lookup k8s.ServiceCatalogLookup,
	lifecycler mode.Lifecycler,
) error {

	// start up the watcher so that events build up on the channel while we're listing events (which happens below)
	watcher, err := iface.Watch(api.ListOptions{LabelSelector: claimLabelSelector})
	if err != nil {
		return err
	}
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
		case evt := <-ch:
			if evt == nil {
				return nil
			}
			logger.Debugf("received event %s", *evt.claim.Claim)
			receiveEvent(ctx, evt, iface, cmNamespacer, lookup, lifecycler)
		case <-ctx.Done():
			logger.Debugf("loop has been stopped")
			return errLoopStopped
		}
	}
}

// determine whether we should process evt (returns true) or skip it (returns false). a few notes:
//
// - consuming MODIFIED events will result in an infinite loop (we'll keep consuming our own events)
// - DELETED events do not indicate that the service should be deleted. transitioning from (evt.claim.Status=="created",evt.claim.Action="create") --> (evt.claim.Status==X,evt.claim.Action=="delete") indicates that
func shouldProcessEvent(evt *Event) bool {
	// process if the event was just added
	if evt.operation == watch.Added {
		return true
	}
	// process if the event has been created but is now being requested to be deleted
	if evt.operation == watch.Modified &&
		evt.claim.Claim.Status == mode.StatusCreated &&
		evt.claim.Claim.Action == mode.ActionDelete {
		return true
	}
	return false
}

func receiveEvent(
	ctx context.Context,
	evt *Event,
	iface Interactor,
	cmNamespacer kcl.ConfigMapsNamespacer,
	lookup k8s.ServiceCatalogLookup,
	lifecycler mode.Lifecycler,
) {

	if !shouldProcessEvent(evt) {
		logger.Debugf("received operation %s for claim %s, skipping", evt.operation, *evt.claim)
		return
	}

	claimCh := make(chan mode.ServicePlanClaim)
	wrapper := evt.claim
	go processEvent(ctx, evt, cmNamespacer, lookup, lifecycler, claimCh)
	stopLoop := false
	for {
		select {
		case claim := <-claimCh:
			wrapper.Claim = &claim
			w, err := iface.Update(wrapper)
			if err != nil {
				logger.Errorf("error updating claim %s (%s), stopping", claim.ClaimID, err)
				stopLoop = true
				break
			}
			if claim.Status == mode.StatusFailed ||
				claim.Status == mode.StatusCreated ||
				claim.Status == mode.StatusDeleted {
				stopLoop = true
				break
			}
			wrapper = w
		case <-ctx.Done():
			return
		}
		if stopLoop {
			break
		}
	}
}
