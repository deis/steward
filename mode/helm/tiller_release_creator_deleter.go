package helm

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type tillerRCD struct {
	tillerHost string
}

func (t tillerRCD) Create(ch *chart.Chart, installNS string) (*rls.InstallReleaseResponse, error) {
	c, err := grpc.Dial(t.tillerHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rlsCl := rls.NewReleaseServiceClient(c)
	ctx := context.Background()
	req := &rls.InstallReleaseRequest{
		Chart:        ch,
		Namespace:    installNS,
		ReuseName:    true,
		DisableHooks: false,
	}
	logger.Debugf("installing release for chart %s", ch.Metadata.Name)
	return rlsCl.InstallRelease(ctx, req)
}

func (t tillerRCD) Delete(relName string) (*rls.UninstallReleaseResponse, error) {
	c, err := grpc.Dial(t.tillerHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rlsCl := rls.NewReleaseServiceClient(c)
	ctx := context.Background()
	req := &rls.UninstallReleaseRequest{
		Name:         relName,
		DisableHooks: false,
	}
	logger.Debugf("uninstalling release %s", relName)
	return rlsCl.UninstallRelease(ctx, req)
}

// NewTillerReleaseCreatorDeleter returns a new ReleaseCreatorDeleter implemented with a tiller backend
func NewTillerReleaseCreatorDeleter(tillerHost string) ReleaseCreatorDeleter {
	return tillerRCD{tillerHost: tillerHost}
}
