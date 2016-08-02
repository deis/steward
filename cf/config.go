package cf

import (
	"github.com/deis/steward/web"
	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "steward"
)

// Config is the envconfig-compatible struct for a backing CF service broker
type Config struct {
	Scheme   string `envconfig:"CF_BROKER_SCHEME" required:"true"`
	Hostname string `envconfig:"CF_BROKER_HOSTNAME" required:"true"`
	Port     int    `envconfig:"CF_BROKER_PORT" required:"true"`
	Username string `envconfig:"CF_BROKER_USERNAME" required:"true"`
	Password string `envconfig:"CF_BROKER_PASSWORD" required:"true"`
}

// GetConfig gets the
func GetConfig() (*Config, error) {
	ret := new(Config)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// BasicAuth returns the basic auth struct that this CF broker config represents
func (c Config) BasicAuth() *web.BasicAuth {
	return &web.BasicAuth{Username: c.Username, Password: c.Password}
}
