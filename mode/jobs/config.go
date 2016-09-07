package jobs

import (
	"time"

	c "github.com/deis/steward/config"
)

type config struct {
	PodNamespace     string `envconfig:"POD_NAMESPACE" required:"true"`
	Image            string `envconfig:"JOBS_IMAGE" required:"true"`
	ConfigMapName    string `envconfig:"JOBS_CONFIG_MAP"`
	SecretName       string `envconfig:"JOBS_SECRET"`
	PollIntervalMSec int    `envconfig:"JOBS_POLL_INTERVAL" default:"10000"` // 10 seconds
	TimeoutMSec      int    `envconfig:"JOBS_TIMEOUT" default:"900000"`      // 15 minutes
}

// getConfig obtains jobs mode configuration from environment variables.
func getConfig() (*config, error) {
	ret := &config{}
	if err := c.Load(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// getPollInterval returns the frequency with which to check the status of a job.
func (c *config) getPollInterval() time.Duration {
	return time.Duration(time.Duration(c.PollIntervalMSec) * time.Millisecond)
}

// getTimeout returns the maximum time to wait for job completion.
func (c *config) getTimeout() time.Duration {
	return time.Duration(time.Duration(c.TimeoutMSec) * time.Millisecond)
}
