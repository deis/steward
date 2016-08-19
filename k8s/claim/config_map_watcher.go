package claim

import (
	"errors"

	"k8s.io/kubernetes/pkg/watch"
)

var (
	errNotAConfigMap = errors.New("not a config map")
)

type configMapWatcher struct {
	watchIface watch.Interface
	closeCh    chan struct{}
}

// NewWatcher returns a watcher that uses watchIface to get and return events
func newConfigMapWatcher(watchIface watch.Interface) Watcher {
	return &configMapWatcher{
		watchIface: watchIface,
		closeCh:    make(chan struct{}),
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
			case rawEvt, open := <-c.watchIface.ResultChan():
				// bail out of the goroutine if the internal result channel was closed
				if !open {
					logger.Debugf("internal watch channel was closed. closing down watch goroutine")
					return
				}
				evt, err := eventFromConfigMapEvent(rawEvt)
				if err != nil {
					logger.Debugf("error converting raw event to service plan claim event (%s). continuing", err)
				} else {
					select {
					case retCh <- evt:
					case <-c.closeCh:
						logger.Debugf("loop was stopped while trying to send a new event, returning")
						return
					}
				}
			}
		}
	}()
	return retCh
}

func (c *configMapWatcher) Stop() {
	c.watchIface.Stop()
	close(c.closeCh)
}
