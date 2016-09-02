package claim

import (
	"context"
	"errors"
	"time"

	"github.com/deis/steward/k8s"
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
		case evt, open := <-ch:
			// if the watch channel was closed, do the following in order:
			//
			// 1. check to see if the loop should be shut down. if so, return immediately
			// 2. otherwise, re-open the watch and continue to the next iteration of the loop
			if !open {
				// this if statement is this is semantically equivalent to if evt == nil {...}
				select {
				case <-ctx.Done():
					logger.Debugf("loop has been stopped")
					return errLoopStopped
				default:
				}
				logger.Debugf("watch channel was closed, pausing then re-opening the watch and continuing")
				time.Sleep(1 * time.Second)
				ch = watcher.ResultChan()
				continue
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

	claimUpdateCh := make(chan claimUpdate)
	wrapper := evt.claim

	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go nextAction(cancelCtx, evt, cmNamespacer, lookup, lifecycler, claimUpdateCh)

	for {
		select {
		case claimUpdate := <-claimUpdateCh:
			// stop watching the processor if it failed
			if claimUpdate.err != nil {
				logger.Errorf("error in claim processing (%s)", claimUpdate.err)
				return
			}

			// update the claim in k8s
			wrapper.Claim = &claimUpdate.newClaim
			w, err := iface.Update(wrapper)
			if err != nil {
				logger.Errorf("error updating claim %s (%s), stopping", claimUpdate.newClaim.ClaimID, err)
				return
			}

			// if stop is true, then the claim was to be updated a final time, and then the loop was to stop listening
			if claimUpdate.stop {
				return
			}
			wrapper = w
		case <-ctx.Done():
			return
		}
	}
}
