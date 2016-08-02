package k8s

import (
	"time"

	"k8s.io/kubernetes/pkg/client/restclient"
	// "k8s.io/kubernetes/pkg/watch"
	"github.com/juju/loggo"
)

// NOTE: all code in this file is not used pending resolution of https://github.com/deis/steward/issues/17

// StartLoop starts an infinite loop that polls third party resources on the given apiEndpoint for the given namespace. It's intended to be called in a goroutine.
//
// This function sends any errors encountered on errCh, and closes errCh if the loop terminates for any reason (including if an error occurred). Callers also must pass a chan struct{} to this function, which can be closed to terminate the loop that this function runs
func StartLoop(logger loggo.Logger, cl *restclient.RESTClient, namespace string, stopCh <-chan struct{}, errCh chan<- error) {
	defer close(errCh)
	claims, err := getServicePlanClaims(cl, namespace)
	if err != nil {
		logger.Debugf("error getting service plan claims (%s)", err)
		errCh <- err
		return
	}
	for _, claim := range claims {
		processClaim(logger, claim)
	}

	// TODO: make the sleep duration configurable
	watchStopCh, watchSPCCh, watchErrCh := watchServicePlanClaims(cl, namespace, 1*time.Second)
	defer close(watchStopCh)
	for {
		select {
		case <-stopCh:
			logger.Debugf("loop stopped for namespace %s", namespace)
			return
		case spc := <-watchSPCCh:
			processClaim(logger, spc)
		case err := <-watchErrCh:
			logger.Debugf("loop got error (%s)", err)
			errCh <- err
			return
		}
	}
}

func processClaim(logger loggo.Logger, claim *ServicePlanClaim) {
	logger.Infof("TODO: processing claim %s", *claim)
}
