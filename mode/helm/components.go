package helm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"google.golang.org/grpc"
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
)

// GetComponents returns tiller-backed implementations of the Cataloger and Lifecycler interfaces. Also returns the underlying gRPC connection used for Tiller communications. Callers should close this connection when done communicating with tiller
func GetComponents(
	ctx context.Context,
	httpCl *http.Client,
	cmNamespacer v1.ConfigMapsGetter,
) (*grpc.ClientConn, mode.Cataloger, *mode.Lifecycler, error) {
	helmCfg, err := getConfig()
	if err != nil {
		logger.Errorf("getting helm config (%s)", err)
		return nil, nil, nil, err
	}
	logger.Infof("starting in Helm mode with Tiller backend at %s:%d", helmCfg.TillerIP, helmCfg.TillerPort)

	provBehavior, err := provisionBehaviorFromString(helmCfg.ProvisionBehavior)
	if err != nil {
		logger.Errorf("parsing provision behavior from '%s' (%s)", helmCfg.ProvisionBehavior, err)
		return nil, nil, nil, err
	}

	chart, _, err := getChart(ctx, httpCl, helmCfg.ChartURL)
	if err != nil {
		logger.Errorf("getting chart from %s (%s)", helmCfg.ChartURL, err)
		return nil, nil, nil, err
	}

	tillerHost := fmt.Sprintf("%s:%d", helmCfg.TillerIP, helmCfg.TillerPort)
	conn, err := grpc.Dial(tillerHost, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil, err
	}

	creatorDeleter := newTillerReleaseCreatorDeleter(conn)

	cataloger := newCataloger(helmCfg)
	lifecycler, err := newLifecycler(ctx, chart, helmCfg.ChartInstallNS, provBehavior, creatorDeleter, cmNamespacer)
	if err != nil {
		logger.Errorf("creating a new helm mode lifecycler (%s)", err)
		conn.Close()
		return nil, nil, nil, err
	}

	return conn, cataloger, lifecycler, nil
}
