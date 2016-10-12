package claim

import (
	"context"

	"github.com/deis/steward/k8s"
	"k8s.io/client-go/1.4/pkg/api"
)

// Interactor is the interface that enables the requisite set of operations on claims that the control loop needs to do its job
type Interactor interface {
	Get(string) (*k8s.ServicePlanClaimWrapper, error)
	List(opts api.ListOptions) (*k8s.ServicePlanClaimsListWrapper, error)
	Update(*k8s.ServicePlanClaimWrapper) (*k8s.ServicePlanClaimWrapper, error)
	Watch(context.Context, api.ListOptions) Watcher
}
