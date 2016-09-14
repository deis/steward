package helm

import (
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type createCall struct {
	chart     *chart.Chart
	installNS string
}

type fakeCreatorDeleter struct {
	createResp    *rls.InstallReleaseResponse
	createRespErr error
	deleteResp    *rls.UninstallReleaseResponse
	deleteRespErr error
	createCalls   []*createCall
	deleteCalls   []string
}

func (f *fakeCreatorDeleter) Create(chart *chart.Chart, installNS string) (*rls.InstallReleaseResponse, error) {
	f.createCalls = append(f.createCalls, &createCall{chart: chart, installNS: installNS})
	return f.createResp, f.createRespErr
}

func (f *fakeCreatorDeleter) Delete(releaseName string) (*rls.UninstallReleaseResponse, error) {
	f.deleteCalls = append(f.deleteCalls, releaseName)
	return f.deleteResp, f.deleteRespErr
}
