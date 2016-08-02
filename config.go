package main

import (
	"fmt"

	"github.com/deis/steward/web"
	"github.com/juju/loggo"
	"github.com/kelseyhightower/envconfig"
)

const (
	cfMode = "cf"
)

type errModeUnsupported struct {
	mode string
}

func (e errModeUnsupported) Error() string {
	return fmt.Sprintf("mode '%s' is unsupported", e.mode)
}

type config struct {
	Mode              string `envconfig:"MODE" default:"cf"`
	LogLevel          string `envconfig:"LOG_LEVEL" default:"info"`
	ServerPort        int    `envconfig:"BROKER_API_SERVER_PORT" default:"8080"`
	BrokerAPIUsername string `envconfig:"BROKER_API_USERNAME" default:"deis"`
	BrokerAPIPassword string `envconfig:"BROKER_API_PASSWORD" default:"steward"`
}

func (c config) validate() error {
	switch c.Mode {
	case cfMode:
	default:
		return errModeUnsupported{mode: c.Mode}
	}
	return nil
}

func (c config) logLevel() loggo.Level {
	switch c.LogLevel {
	case "trace":
		return loggo.TRACE
	case "debug":
		return loggo.DEBUG
	case "info":
		return loggo.INFO
	case "warning":
		return loggo.WARNING
	case "error":
		return loggo.ERROR
	case "critical":
		return loggo.CRITICAL
	default:
		return loggo.INFO
	}
}

func (c config) hostString() string {
	return fmt.Sprintf(":%d", c.ServerPort)
}

func (c config) basicAuth() *web.BasicAuth {
	return &web.BasicAuth{Username: c.BrokerAPIUsername, Password: c.BrokerAPIPassword}
}

func getConfig(appName string) (*config, error) {
	spec := new(config)
	if err := envconfig.Process(appName, spec); err != nil {
		return nil, err
	}
	return spec, nil
}
