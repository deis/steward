package main

import (
	"os"

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
	logger.SetLogLevel(loggo.TRACE)
	logger.Infof("steward version %s started", version)
	cfg, err := getConfig(appName)
	if err != nil {
		logger.Criticalf("error getting config (%s)", err)
		os.Exit(1)
	}
	logger.SetLogLevel(cfg.logLevel())
	if err := cfg.validate(); err != nil {
		logger.Criticalf("error with config (%s)", err)
		os.Exit(1)
	}

	k8sClient, err := kcl.NewInCluster()
	if err != nil {
		logger.Criticalf("error creating new k8s client (%s)", err)
		os.Exit(1)
	}

	errCh := make(chan error)

	switch cfg.Mode {
	case cfMode:
		if err := runCFMode(
			logger,
			cfg.hostString(),
			cfg.basicAuth(),
			k8sClient.RESTClient,
			errCh,
			k8s.NewConfigMapCreator(k8sClient),
			k8s.NewSecretCreator(k8sClient),
		); err != nil {
			logger.Criticalf("error executing in CloudFoundry mode (%s)", err)
			os.Exit(1)
		}
	default:
		logger.Criticalf("no catalog to publish for mode %s", cfg.Mode)
		os.Exit(1)
	}

	// NOTE: this code is pending resolution of https://github.com/deis/steward/issues/17
	// namespaces := []string{"default", "deis", "steward"}
	// go func() {
	// 	k8s.StartLoops(logger, k8sClient.RESTClient, namespaces, stopCh, errCh)
	// }()

	// TODO: listen for signal and delete all service catalog entries before quitting
	select {
	case err := <-errCh:
		if err != nil {
			logger.Criticalf("%s", err)
			os.Exit(1)
		} else {
			logger.Criticalf("unknown error, crashing")
			os.Exit(1)
		}
	}
}
