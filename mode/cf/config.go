package cf

import (
	"github.com/deis/steward/config"
	"github.com/deis/steward/web"
)

// Config is the envconfig-compatible struct for a backing CloudFoundry broker
type Config struct {
	Scheme   string `envconfig:"CF_BROKER_SCHEME" required:"true"`
	Hostname string `envconfig:"CF_BROKER_HOSTNAME" required:"true"`
	Port     int    `envconfig:"CF_BROKER_PORT" required:"true"`
	Username string `envconfig:"CF_BROKER_USERNAME" required:"true"`
	Password string `envconfig:"CF_BROKER_PASSWORD" required:"true"`
}

// GetConfig gets the configuration for CF mode
func GetConfig() (*Config, error) {
	ret := new(Config)
	if err := config.Load(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c Config) basicAuth() *web.BasicAuth {
	return &web.BasicAuth{Username: c.Username, Password: c.Password}
}
