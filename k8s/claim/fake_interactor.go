package claim

import (
	"context"

	"k8s.io/client-go/1.4/pkg/api"
)

// FakeInteractor is a fake implementation of Interactor, to be use in unit testing
type FakeInteractor struct {
}

// Get is the Interactor interface implementation
func (f *FakeInteractor) Get(string) (*ServicePlanClaimWrapper, error) {
	return nil, nil
}

// List is the Interactor interface implementation
func (f *FakeInteractor) List(opts api.ListOptions) (*ServicePlanClaimsListWrapper, error) {
	return nil, nil
}

// Update is the Interactor interface implementation
func (f *FakeInteractor) Update(*ServicePlanClaimWrapper) (*ServicePlanClaimWrapper, error) {
	return nil, nil
}

// Watch is the Interactor interface implementation
func (f *FakeInteractor) Watch(context.Context, api.ListOptions) Watcher {
	return nil
}
