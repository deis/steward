package k8s

import (
	"github.com/juju/loggo"
	"k8s.io/kubernetes/pkg/client/restclient"
)

// NOTE: all code in this file is not used pending resolution of https://github.com/deis/steward/issues/17

// StartLoops calls StartLoop in a goroutine for each namespace in namespaces. For each call to StartLoop, it calls createTPRI to get a new ThirdPartyResourceInterface to pass into StartLoop.
//
// This function should be called inside a new goroutine. If it encounters an error in any StartLoop goroutine, it sends the error on errCh, shuts down all StartLoop goroutines, and closes errCh. If stopCh is closed by the caller, all StartLoop goroutines will be shut down, but errCh will not be closed.
func StartLoops(
	logger loggo.Logger,
	cl *restclient.RESTClient,
	namespaces []string,
	stopCh <-chan struct{},
	errCh chan<- error,
) {

	internalStopCh := make(chan struct{})
	internalErrCh := make(chan error)
	for _, ns := range namespaces {
		go func(ns string) {
			StartLoop(logger, cl, ns, internalStopCh, internalErrCh)
		}(ns)
	}
	select {
	case err := <-internalErrCh:
		errCh <- err
		close(internalStopCh)
	case <-stopCh:
		close(internalStopCh)
	}
}
