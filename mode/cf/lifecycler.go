package cf

import (
	"github.com/deis/steward/mode"
)

// NewLifecycler returns a new mode.Lifecycler that's implemented with a backend CF broker
func NewLifecycler(cl *RESTClient) *mode.Lifecycler {
	return &mode.Lifecycler{
		Provisioner:   NewProvisioner(cl),
		Deprovisioner: NewDeprovisioner(cl),
		Binder:        NewBinder(cl),
		Unbinder:      NewUnbinder(cl),
	}
}
