package main

import (
	"os"

	modeutils "github.com/deis/steward/mode/utils"
	"github.com/deis/steward/web/api"
	"github.com/juju/loggo"
)

const (
	appName = "steward"
)

var (
	logger     = loggo.GetLogger("")
	version    = "dev"
	namespaces = []string{"steward", "default", "deis"}
)

func main() {
	logger.Infof("steward version %s started", version)
	cfg, err := getConfig(appName)
	if err != nil {
		logger.Criticalf("error getting config (%s)", err)
		os.Exit(1)
	}
	logger.SetLogLevel(cfg.logLevel())

	errCh := make(chan error)

	if err := modeutils.Run(cfg.Mode, errCh, namespaces); err != nil {
		logger.Criticalf("Error starting %s mode: %s", cfg.Mode, err)
		os.Exit(1)
	}

	// NOTE: this code is pending resolution of https://github.com/deis/steward/issues/17
	// namespaces := []string{"default", "deis", "steward"}
	// go func() {
	// 	k8s.StartLoops(k8sClient.RESTClient, namespaces, stopCh, errCh)
	// }()

	// Start the API server
	go api.Serve(errCh)

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
