package claim

import (
	"context"
	"errors"

	"k8s.io/client-go/1.4/pkg/watch"
)

var (
	errNotAConfigMap = errors.New("not a config map")
	errStopped       = errors.New("stopped")
)

type configMapWatcher struct {
	ctx     context.Context
	ifaceFn func() (watch.Interface, error)
}

// NewWatcher returns a watcher that uses watchIface to get and return events
func newConfigMapWatcher(ctx context.Context, ifaceFn func() (watch.Interface, error)) Watcher {
	return &configMapWatcher{
		ctx:     ctx,
		ifaceFn: ifaceFn,
	}
}

// receives on iface.ResultChan() until either that channel or closeCh was closed. returned errWatchClosed in the former case, errStopped in the latter
func watchLoop(ctx context.Context, iface watch.Interface, retCh chan<- *Event) error {
	defer iface.Stop()
	resCh := iface.ResultChan()
	for {
		select {
		case <-ctx.Done():
			return errStopped
		case rawEvt, open := <-resCh:
			if !open {
				return errWatchClosed
			}
			evt, err := eventFromConfigMapEvent(rawEvt)
			if err != nil {
				logger.Debugf("error converting raw event to service plan claim event (%s). continuing", err)
			} else {
				select {
				case retCh <- evt:
				case <-ctx.Done():
					logger.Debugf("loop was stopped while trying to send a new event, returning")
					return errStopped
				}
			}
		}
	}
}

// ResultChan is the (k8s.io/kubernetes/pkg/watch).Interface interface implementation. It returns a channel that will be closed either when Stop() is called, or when the server severs the connection, which may happen intermittently.
func (c *configMapWatcher) ResultChan() <-chan *Event {
	retCh := make(chan *Event)
	go func() {
		defer close(retCh)
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
			}
			iface, err := c.ifaceFn()
			if err != nil {
				logger.Errorf("error getting a new watch interface (%s)", err)
				return
			}
			watchErr := watchLoop(c.ctx, iface, retCh)
			if watchErr == errWatchClosed {
				logger.Debugf("Kubernetes watch API closed connection. Re-opening")
				continue
			} else {
				return
			}
		}
	}()
	return retCh
}
