package helm

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type tillerRCD struct {
	conn *grpc.ClientConn
}

func (t tillerRCD) Create(ch *chart.Chart, installNS string) (*rls.InstallReleaseResponse, error) {
	rlsCl := rls.NewReleaseServiceClient(t.conn)
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
	rlsCl := rls.NewReleaseServiceClient(t.conn)
	ctx := context.Background()
	req := &rls.UninstallReleaseRequest{
		Name:         relName,
		DisableHooks: false,
	}
	logger.Debugf("uninstalling release %s", relName)
	return rlsCl.UninstallRelease(ctx, req)
}

// newTillerReleaseCreatorDeleter returns a new ReleaseCreatorDeleter implemented with a tiller backend
func newTillerReleaseCreatorDeleter(conn *grpc.ClientConn) ReleaseCreatorDeleter {
	return tillerRCD{conn: conn}
}
