package main

import (
	"log"

	"github.com/deis/steward/k8s"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

var version = "dev"

func main() {
	log.Printf("steward version %s started", version)
	k8sClient, err := kcl.NewInCluster()
	if err != nil {
		log.Fatalf("error creating new k8s client (%s)", err)
	}
	errCh := make(chan error)
	stopCh := make(chan struct{})
	namespaces := []string{"default", "deis"}
	go func() {
		k8s.StartLoops(k8sClient.RESTClient, namespaces, stopCh, errCh)
	}()
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("error running control loop for 'deis' (%s)", err)
		} else {
			log.Fatalf("control loop for 'deis' stopped for unknown reason")
		}
	}
}
