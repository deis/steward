package helm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
)

// GetComponents returns suitable implementations of the Cataloger and Lifecycler interfaces
func GetComponents(
	ctx context.Context,
	httpCl *http.Client,
	cmNamespacer v1.ConfigMapsGetter,
) (mode.Cataloger, *mode.Lifecycler, error) {
	helmCfg, err := getConfig()
	if err != nil {
		logger.Errorf("getting helm config (%s)", err)
		return nil, nil, err
	}
	logger.Infof("starting in Helm mode with Tiller backend at %s:%d", helmCfg.TillerIP, helmCfg.TillerPort)

	provBehavior, err := provisionBehaviorFromString(helmCfg.ProvisionBehavior)
	if err != nil {
		logger.Errorf("parsing provision behavior from '%s' (%s)", helmCfg.ProvisionBehavior, err)
		return nil, nil, err
	}

	chart, _, err := getChart(ctx, httpCl, helmCfg.ChartURL)
	if err != nil {
		logger.Errorf("getting chart from %s (%s)", helmCfg.ChartURL, err)
		return nil, nil, err
	}

	tillerHost := fmt.Sprintf("%s:%d", helmCfg.TillerIP, helmCfg.TillerPort)
	creatorDeleter := newTillerReleaseCreatorDeleter(tillerHost)

	cataloger := newCataloger(helmCfg)
	lifecycler, err := newLifecycler(ctx, chart, helmCfg.ChartInstallNS, provBehavior, creatorDeleter, cmNamespacer)
	if err != nil {
		logger.Errorf("creating a new helm mode lifecycler (%s)", err)
		return nil, nil, err
	}

	return cataloger, lifecycler, nil
}
