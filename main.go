package main

import (
	"context"
	"net/http"
	"os"

	modeutils "github.com/deis/steward/mode/utils"
	"github.com/deis/steward/web/api"
	"github.com/juju/loggo"
)

var (
	logger  = loggo.GetLogger("")
	version = "dev"
)

func exitWithCode(cancelFn func(), exitCode int) {
	cancelFn()
	os.Exit(exitCode)
}

func main() {
	logger.Infof("steward version %s started", version)
	cfg, err := getRootConfig()
	if err != nil {
		logger.Criticalf("error getting config (%s)", err)
		os.Exit(1)
	}
	logger.SetLogLevel(cfg.logLevel())

	errCh := make(chan error)
	rootCtx := context.Background()
	httpCl := http.DefaultClient
	ctx, cancelFn := context.WithCancel(rootCtx)
	defer cancelFn()

	if err := modeutils.Run(ctx, httpCl, cfg.Mode, errCh, cfg.WatchNamespaces); err != nil {
		logger.Criticalf("Error starting %s mode: %s", cfg.Mode, err)
		exitWithCode(cancelFn, 1)
	}

	// Start the API server
	go api.Serve(errCh)

	// TODO: listen for signal and delete all service catalog entries before quitting
	select {
	case err := <-errCh:
		if err != nil {
			logger.Criticalf("%s", err)
			exitWithCode(cancelFn, 1)
		} else {
			logger.Criticalf("unknown error, crashing")
			exitWithCode(cancelFn, 1)
		}
	}
}
