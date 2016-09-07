package claim

import (
	"errors"

	"k8s.io/kubernetes/pkg/watch"
)

var (
	errNotAConfigMap = errors.New("not a config map")
	errStopped       = errors.New("stopped")
)

type configMapWatcher struct {
	ifaceFn func() (watch.Interface, error)
	closeCh chan struct{}
}

// NewWatcher returns a watcher that uses watchIface to get and return events
func newConfigMapWatcher(ifaceFn func() (watch.Interface, error)) Watcher {
	return &configMapWatcher{
		ifaceFn: ifaceFn,
		closeCh: make(chan struct{}),
	}
}

// receives on iface.ResultChan() until either that channel or closeCh was closed. returned errWatchClosed in the former case, errStopped in the latter
func watchLoop(iface watch.Interface, retCh chan<- *Event, closeCh <-chan struct{}) error {
	defer iface.Stop()
	resCh := iface.ResultChan()
	for {
		select {
		case <-closeCh:
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
				case <-closeCh:
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
			case <-c.closeCh:
				return
			default:
			}
			iface, err := c.ifaceFn()
			if err != nil {
				logger.Errorf("error getting a new watch interface (%s)", err)
				return
			}
			watchErr := watchLoop(iface, retCh, c.closeCh)
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

func (c *configMapWatcher) Stop() {
	close(c.closeCh)
}
