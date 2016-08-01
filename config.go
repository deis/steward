package main

import (
	"fmt"

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

type errMissingCFHostNameOrPassword struct {
	hostname string
	password string
}

func (e errMissingCFHostNameOrPassword) Error() string {
	if e.hostname == "" && e.password == "" {
		return fmt.Sprintf("missing CF hostname and password")
	}
	if e.hostname == "" {
		return fmt.Sprintf("missing CF hostname")
	}
	if e.password == "" {
		return fmt.Sprintf("missing CF password")
	}
	// we should never get here
	return fmt.Sprintf("all CF config parameters are present. you should never see this error, because it is not an error!")
}

type config struct {
	Mode             string `envconfig:"STEWARD_MODE" default:"cf"`
	CFBrokerScheme   string `envconfig:"STEWARD_CF_BROKER_SCHEME" default:"http"`
	CFBrokerHostname string `envconfig:"STEWARD_CF_BROKER_HOSTNAME" default:"localhost:8080"`
	CFBrokerUsername string `envconfig:"STEWARD_CF_BROKER_USERNAME" default:"admin"`
	CFBrokerPassword string `envconfig:"STEWARD_CF_BROKER_PASSWORD" default:"password"`
	LogLevel         string `envconfig:"STEWARD_LOG_LEVEL" default:"info"`
}

func (c config) validate() error {
	switch c.Mode {
	case cfMode:
		if c.CFBrokerHostname == "" || c.CFBrokerPassword == "" {
			return errMissingCFHostNameOrPassword{hostname: c.CFBrokerHostname, password: c.CFBrokerPassword}
		}
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

func getConfig(appName string) (*config, error) {
	spec := new(config)
	if err := envconfig.Process(appName, spec); err != nil {
		return nil, err
	}
	return spec, nil
}
