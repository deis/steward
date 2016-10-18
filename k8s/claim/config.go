package claim

import (
	"time"

	c "github.com/deis/steward/config"
)

type config struct {
	MaxAsyncMinutes int `envconfig:"MAX_ASYNC_MINUTES" default:"60"`
}

// getConfig obtains claim configuration from environment variables
func getConfig() (*config, error) {
	ret := &config{}
	if err := c.Load(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// getMaxAsyncDuration returns the maximum time spent on an aysnchronous provisioning or
// deprovisioning before giving up on polling for the last operation
func (c *config) getMaxAsyncDuration() time.Duration {
	return time.Duration(time.Duration(c.MaxAsyncMinutes) * time.Minute)
}
