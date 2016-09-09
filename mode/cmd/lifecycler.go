package cmd

import (
	"github.com/deis/steward/mode"
)

// newLifecycler returns a new mode.Lifecycler that's implemented with a backend cmd broker
func newLifecycler(pr *podRunner) *mode.Lifecycler {
	return &mode.Lifecycler{
		Provisioner:   newProvisioner(pr),
		Deprovisioner: newDeprovisioner(pr),
		Binder:        newBinder(pr),
		Unbinder:      newUnbinder(pr),
	}
}
