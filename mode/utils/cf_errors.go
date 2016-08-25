package utils

import "fmt"

type errGettingCFBrokerConfig struct {
	Original error
}

func (e errGettingCFBrokerConfig) Error() string {
	return fmt.Sprintf("error getting Cloud Foundry broker config: %s", e.Original)
}
