package k8s

import (
	"log"
	"time"

	"k8s.io/kubernetes/pkg/client/restclient"
	// "k8s.io/kubernetes/pkg/watch"
)

// StartLoop starts an infinite loop that polls third party resources on the given apiEndpoint for the given namespace. It's intended to be called in a goroutine.
//
// This function sends any errors encountered on errCh, and closes errCh if the loop terminates for any reason (including if an error occurred). Callers also must pass a chan struct{} to this function, which can be closed to terminate the loop that this function runs
func StartLoop(cl *restclient.RESTClient, namespace string, stopCh <-chan struct{}, errCh chan<- error) {
	defer close(errCh)
	claims, err := getServicePlanClaims(cl, namespace)
	if err != nil {
		log.Printf("error getting service plan claims (%s)", err)
		errCh <- err
		return
	}
	for _, claim := range claims.Items {
		processClaim(claim)
	}

	// TODO: make the sleep duration configurable
	watchStopCh, watchSPCCh, watchErrCh := watchServicePlanClaims(cl, namespace, 1*time.Second)
	defer close(watchStopCh)
	for {
		select {
		case <-stopCh:
			log.Printf("loop stopped for namespace %s", namespace)
			return
		case spc := <-watchSPCCh:
			processClaim(spc)
		case err := <-watchErrCh:
			log.Printf("got error %s", err)
			errCh <- err
			return
		}
	}
}

func processClaim(claim *ServicePlanClaim) {
	log.Printf("processing claim %s", *claim)
}
