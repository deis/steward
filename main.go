package main

import (
	"net/http"
	"os"

	"github.com/deis/steward/cf"
	"github.com/deis/steward/k8s"
	"github.com/juju/loggo"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	appName = "steward"
)

var (
	version = "dev"
)

func main() {
	logger := loggo.GetLogger("")
	logger.ApplyConfig(loggo.LoggerConfig{Level: loggo.TRACE})
	logger.Infof("steward version %s started", version)
	cfg, err := getConfig(appName)
	if err != nil {
		logger.Criticalf("error getting config (%s)", err)
		os.Exit(1)
	}
	logger.ApplyConfig(loggo.LoggerConfig{Level: cfg.logLevel()})
	if err := cfg.validate(); err != nil {
		logger.Criticalf("error with config (%s)", err)
		os.Exit(1)
	}

	k8sClient, err := kcl.NewInCluster()
	if err != nil {
		logger.Criticalf("error creating new k8s client (%s)", err)
		os.Exit(1)
	}

	switch cfg.Mode {
	case cfMode:
		logger.Infof(
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
		published, err := publishCloudFoundryCatalog(logger, cfClient, k8sClient.RESTClient)
		if err != nil {
			logger.Criticalf("error publishing the cloud foundry service catalog (%s)", err)
			os.Exit(1)
		}
		logger.Infof("published %d entries into the catalog", len(published))
		for _, pub := range published {
			logger.Debugf("%s", pub.Info.Name)
		}
	default:
		logger.Criticalf("no catalog to publish for mode %s", cfg.Mode)
		os.Exit(1)
	}

	errCh := make(chan error)
	stopCh := make(chan struct{})
	namespaces := []string{"default", "deis", "steward"}
	go func() {
		k8s.StartLoops(logger, k8sClient.RESTClient, namespaces, stopCh, errCh)
	}()

	// TODO: listen for signal and delete all service catalog entries before quitting
	select {
	case err := <-errCh:
		if err != nil {
			logger.Criticalf("error running control loop for 'deis' (%s)", err)
			os.Exit(1)
		} else {
			logger.Criticalf("control loop for 'deis' stopped for unknown reason")
			os.Exit(1)
		}
	}
}
