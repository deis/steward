package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	// AppName is the standard app name to use in fetching configs
	AppName = "steward"
)

// Load is a convenience function for calling envconfig.Process(AppName, ret) (godoc.org/github.com/kelseyhightower/envconfig#Process)
func Load(ret interface{}) error {
	return envconfig.Process(AppName, ret)
}
