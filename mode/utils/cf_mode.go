package utils

import (
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
)

func getCfModeComponents() (mode.Cataloger, *mode.Lifecycler, error) {
	cfCfg, err := getCfConfig()
	if err != nil {
		return nil, nil, errGettingBrokerConfig{Original: err}
	}
	logger.Infof(
		"starting in Cloud Foundry mode with hostname %s, port %d, and username %s",
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
	)
	cfClient := cf.NewRESTClient(
		http.DefaultClient,
		cfCfg.Scheme,
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
		cfCfg.Password,
	)
	cataloger := cf.NewCataloger(cfClient)
	lifecycler := cf.NewLifecycler(cfClient)
	return cataloger, lifecycler, nil
}
