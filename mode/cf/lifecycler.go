package cf

import (
	"context"
	"time"

	"github.com/deis/steward/mode"
)

// newLifecycler returns a new mode.Lifecycler that's implemented with a backend CF broker
func newLifecycler(ctx context.Context, cl *restClient, callTimeout time.Duration) *mode.Lifecycler {
	return &mode.Lifecycler{
		Provisioner:         newProvisioner(ctx, cl, callTimeout),
		Deprovisioner:       newDeprovisioner(ctx, cl, callTimeout),
		Binder:              newBinder(ctx, cl, callTimeout),
		Unbinder:            newUnbinder(ctx, cl, callTimeout),
		LastOperationGetter: newLastOperationGetter(ctx, cl, callTimeout),
	}
}
