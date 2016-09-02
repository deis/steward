package cf

import (
	"context"
	"time"

	"github.com/deis/steward/mode"
)

// NewLifecycler returns a new mode.Lifecycler that's implemented with a backend CF broker
func NewLifecycler(ctx context.Context, cl *RESTClient, callTimeout time.Duration) *mode.Lifecycler {
	return &mode.Lifecycler{
		Provisioner:   NewProvisioner(ctx, cl, callTimeout),
		Deprovisioner: NewDeprovisioner(ctx, cl, callTimeout),
		Binder:        NewBinder(ctx, cl, callTimeout),
		Unbinder:      NewUnbinder(ctx, cl, callTimeout),
	}
}
