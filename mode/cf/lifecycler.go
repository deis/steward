package cf

import (
	"github.com/deis/steward/mode"
)

type lifecycler struct {
	mode.Provisioner
	mode.Deprovisioner
	mode.Binder
	mode.Unbinder
}

// NewLifecycler returns a new mode.Lifecycler that's implemented with a backend CF broker
func NewLifecycler(cl *RESTClient) mode.Lifecycler {
	return &lifecycler{
		Provisioner:   NewProvisioner(cl),
		Deprovisioner: NewDeprovisioner(cl),
		Binder:        NewBinder(cl),
		Unbinder:      NewUnbinder(cl),
	}
}
