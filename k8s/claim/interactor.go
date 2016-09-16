package claim

import (
	"context"

	"k8s.io/client-go/1.4/pkg/api"
)

// Interactor is the interface that enables the requisite set of operations on claims that the control loop needs to do its job
type Interactor interface {
	Get(string) (*ServicePlanClaimWrapper, error)
	List(opts api.ListOptions) (*ServicePlanClaimsListWrapper, error)
	Update(*ServicePlanClaimWrapper) (*ServicePlanClaimWrapper, error)
	Watch(context.Context, api.ListOptions) Watcher
}
