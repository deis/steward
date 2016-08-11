package main

import (
	"net/http"
	"os"

	"github.com/deis/steward/mode/cf"
	"github.com/juju/loggo"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	APIScheme       string `envconfig:"API_SCHEME" default:"http"`
	APIHost         string `envconfig:"API_HOST" required:"true"`
	APIPort         int    `envconfig:"API_PORT" required:"true"`
	APIUser         string `envconfig:"API_USER" required:"true"`
	APIPass         string `envconfig:"API_PASS" required:"true"`
	TargetNamespace string `envconfig:"TARGET_NAMESPACE" default:"steward"`
}

const (
	appName = "steward"
)

func main() {
	logger := loggo.GetLogger(appName)
	logger.SetLogLevel(loggo.TRACE)
	cfg := new(config)
	if err := envconfig.Process(appName, cfg); err != nil {
		logger.Criticalf("config error (%s)", err)
		os.Exit(1)
	}
	cl := cf.NewRESTClient(http.DefaultClient, cfg.APIScheme, cfg.APIHost, cfg.APIPort, cfg.APIUser, cfg.APIPass)
	if err := drive(logger, cl, cfg.TargetNamespace); err != nil {
		logger.Criticalf("integration test error (%s)", err)
		os.Exit(1)
	}
}
