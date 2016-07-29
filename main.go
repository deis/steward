package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/deis/steward/cf"
	"github.com/deis/steward/k8s"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	appName = "steward"
)

var (
	version = "dev"
)

func main() {
	log.Printf("steward version %s started", version)
	flag.Parse()
	cfg, err := getConfig(appName)
	if err != nil {
		log.Fatalf("error getting config (%s)", err)
	}
	if err := cfg.validate(); err != nil {
		log.Fatalf("error with config (%s)", err)
	}

	k8sClient, err := kcl.NewInCluster()
	if err != nil {
		log.Fatalf("error creating new k8s client (%s)", err)
	}

	switch cfg.Mode {
	case cfMode:
		log.Printf(
			"starting in CloudFoundry mode with hostname %s and username %s",
			cfg.CFBrokerHostname,
			cfg.CFBrokerUsername,
		)
		cfClient := cf.NewClient(
			http.DefaultClient,
			cfg.CFBrokerScheme,
			cfg.CFBrokerHostname,
			cfg.CFBrokerUsername,
			cfg.CFBrokerPassword,
		)
		published, err := publishCloudFoundryCatalog(cfClient, k8sClient.RESTClient)
		if err != nil {
			log.Fatalf("error publishing the cloud foundry service catalog (%s)", err)
		}
		log.Printf("published %d entries into the catalog:", len(published))
		for _, pub := range published {
			log.Printf("%s", pub.Info.Name)
		}
	default:
		log.Fatalf("no catalog to publish for mode %s", cfg.Mode)
	}

	errCh := make(chan error)
	stopCh := make(chan struct{})
	namespaces := []string{"default", "deis", "steward"}
	go func() {
		k8s.StartLoops(k8sClient.RESTClient, namespaces, stopCh, errCh)
	}()

	// TODO: listen for signal and delete all service catalog entries before quitting
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("error running control loop for 'deis' (%s)", err)
		} else {
			log.Fatalf("control loop for 'deis' stopped for unknown reason")
		}
	}
}
