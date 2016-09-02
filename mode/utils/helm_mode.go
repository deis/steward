package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/helm"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func getHelmModeComponents(
	ctx context.Context,
	httpCl *http.Client,
	cmNamespacer kcl.ConfigMapsNamespacer,
) (mode.Cataloger, *mode.Lifecycler, error) {
	helmCfg, err := helm.GetConfig()
	if err != nil {
		logger.Errorf("getting helm config (%s)", err)
		return nil, nil, errGettingBrokerConfig{Original: err}
	}
	logger.Infof("starting in Helm mode with Tiller backend at %s:%d", helmCfg.TillerIP, helmCfg.TillerPort)

	provBehavior, err := helm.ProvisionBehaviorFromString(helmCfg.ProvisionBehavior)
	if err != nil {
		logger.Errorf("parsing provision behavior from '%s' (%s)", helmCfg.ProvisionBehavior, err)
		return nil, nil, err
	}

	chart, _, err := helm.GetChart(ctx, httpCl, helmCfg.ChartURL)
	if err != nil {
		logger.Errorf("getting chart from %s (%s)", helmCfg.ChartURL, err)
		return nil, nil, err
	}

	tillerHost := fmt.Sprintf("%s:%d", helmCfg.TillerIP, helmCfg.TillerPort)
	creatorDeleter := helm.NewTillerReleaseCreatorDeleter(tillerHost)

	cataloger := helm.NewCataloger(helmCfg)
	lifecycler, err := helm.NewLifecycler(ctx, chart, helmCfg.ChartInstallNS, provBehavior, creatorDeleter, cmNamespacer)
	if err != nil {
		logger.Errorf("creating a new helm mode lifecycler (%s)", err)
		return nil, nil, err
	}

	return cataloger, lifecycler, nil
}
