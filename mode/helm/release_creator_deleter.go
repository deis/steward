package helm

import (
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

// ReleaseCreator is the interface for creating a helm release. It's intended for use in function params for easy mocking
type ReleaseCreator interface {
	Create(*chart.Chart, string) (*rls.InstallReleaseResponse, error)
}

// ReleaseDeleter is the interface for deleting a helm release. It's intended for use in function params for easy mocking
type ReleaseDeleter interface {
	Delete(string) (*rls.UninstallReleaseResponse, error)
}

// ReleaseCreatorDeleter is the concrete composition of a ReleaseCreator and ReleaseDeleter
type ReleaseCreatorDeleter interface {
	ReleaseCreator
	ReleaseDeleter
}
