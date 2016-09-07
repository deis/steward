package cf

import (
	"context"
	"net/http"

	"github.com/deis/steward/mode"
)

// GetComponents returns suitable implementations of the Cataloger and Lifecycler interfaces
func GetComponents(ctx context.Context, httpCl *http.Client) (mode.Cataloger, *mode.Lifecycler, error) {
	cfCfg, err := getConfig()
	if err != nil {
		return nil, nil, err
	}
	logger.Infof(
		"starting in Cloud Foundry mode with hostname %s, port %d, and username %s",
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
	)
	cfClient := newRESTClient(
		httpCl,
		cfCfg.Scheme,
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
		cfCfg.Password,
	)
	callTimeout := cfCfg.HttpRequestTimeoutSec()
	cataloger := newCataloger(ctx, cfClient, callTimeout)
	lifecycler := newLifecycler(ctx, cfClient, callTimeout)
	return cataloger, lifecycler, nil
}
