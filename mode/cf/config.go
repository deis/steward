package cf

import (
	"time"

	conf "github.com/deis/steward/config"
	"github.com/deis/steward/web"
)

// config is the envconfig-compatible struct for a backing CloudFoundry broker
type config struct {
	Scheme                string `envconfig:"CF_BROKER_SCHEME" required:"true"`
	Hostname              string `envconfig:"CF_BROKER_HOSTNAME" required:"true"`
	Port                  int    `envconfig:"CF_BROKER_PORT" required:"true"`
	Username              string `envconfig:"CF_BROKER_USERNAME" required:"true"`
	Password              string `envconfig:"CF_BROKER_PASSWORD" required:"true"`
	HTTPRequestTimeoutSec int    `envconfig:"HTTP_REQUEST_TIMEOUT_SEC" default:"5"`
}

// getConfig gets the configuration for CF mode
func getConfig() (*config, error) {
	ret := new(config)
	if err := conf.Load(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c config) basicAuth() *web.BasicAuth {
	return &web.BasicAuth{Username: c.Username, Password: c.Password}
}

// HTTPRequestTimeoutSec returns the HTTP request timeout defined on c, as a time.Duration
func (c config) HttpRequestTimeoutSec() time.Duration {
	return time.Duration(c.HTTPRequestTimeoutSec) * time.Second
}
