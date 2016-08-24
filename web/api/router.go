package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var cfg *config
var router *mux.Router

func init() {
	var err error
	cfg, err = parseConfig()
	if err != nil {
		logger.Criticalf("error getting config (%s)", err)
		os.Exit(1)
	}

	router = mux.NewRouter()
	router.StrictSlash(true)

	// healthz
	router.Handle("/healthz", http.HandlerFunc(healthzHandler)).Methods("GET")
}

// Serve starts an HTTP server that handles all inbound requests. This function blocks while the
// server runs, so it should be run in its own goroutine.
func Serve(errCh chan<- error) {
	logger.Infof("starting API server on port %d", cfg.Port)
	host := fmt.Sprintf(":%d", cfg.Port)
	if err := http.ListenAndServe(host, router); err != nil {
		errCh <- err
	}
}
