package utils

import (
	"github.com/deis/steward/web"
	"github.com/kelseyhightower/envconfig"
)

type cfConfig struct {
	Scheme   string `envconfig:"CF_BROKER_SCHEME" required:"true"`
	Hostname string `envconfig:"CF_BROKER_HOSTNAME" required:"true"`
	Port     int    `envconfig:"CF_BROKER_PORT" required:"true"`
	Username string `envconfig:"CF_BROKER_USERNAME" required:"true"`
	Password string `envconfig:"CF_BROKER_PASSWORD" required:"true"`
}

func getCfConfig() (*cfConfig, error) {
	ret := new(cfConfig)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c cfConfig) basicAuth() *web.BasicAuth {
	return &web.BasicAuth{Username: c.Username, Password: c.Password}
}
