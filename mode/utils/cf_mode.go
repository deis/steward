package utils

import (
	"context"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
)

func getCfModeComponents(ctx context.Context, httpCl *http.Client) (mode.Cataloger, *mode.Lifecycler, error) {
	cfCfg, err := cf.GetConfig()
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
		httpCl,
		cfCfg.Scheme,
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
		cfCfg.Password,
	)
	callTimeout := cfCfg.HttpRequestTimeoutSec()
	cataloger := cf.NewCataloger(ctx, cfClient, callTimeout)
	lifecycler := cf.NewLifecycler(ctx, cfClient, callTimeout)
	return cataloger, lifecycler, nil
}
