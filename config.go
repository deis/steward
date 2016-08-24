package main

import (
	"fmt"

	"github.com/juju/loggo"
	"github.com/kelseyhightower/envconfig"
)

type errModeUnsupported struct {
	mode string
}

func (e errModeUnsupported) Error() string {
	return fmt.Sprintf("mode '%s' is unsupported", e.mode)
}

type config struct {
	Mode     string `envconfig:"MODE" default:"cf"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
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

func getConfig(appName string) (*config, error) {
	spec := new(config)
	if err := envconfig.Process(appName, spec); err != nil {
		return nil, err
	}
	return spec, nil
}
