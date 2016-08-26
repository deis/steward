package main

import (
	"github.com/deis/steward/config"
	"github.com/juju/loggo"
)

type rootConfig struct {
	Mode            string   `envconfig:"MODE" default:"cf"`
	LogLevel        string   `envconfig:"LOG_LEVEL" default:"info"`
	WatchNamespaces []string `envconfig:"WATCH_NAMESPACES" default:"default"`
}

func getRootConfig() (*rootConfig, error) {
	ret := new(rootConfig)
	if err := config.Load(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c rootConfig) logLevel() loggo.Level {
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
