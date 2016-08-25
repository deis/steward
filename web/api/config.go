package api

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "steward"
)

type config struct {
	Port int `envconfig:"API_PORT" default:"8080"`
}

func parseConfig() (*config, error) {
	ret := new(config)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
