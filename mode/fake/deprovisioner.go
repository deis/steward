package fake

import (
	"github.com/deis/steward/mode"
)

// Deprovisioner is a fake implementation of (github.com/deis/steward/mode).Deprovisioner, suitable for usage in unit tests
type Deprovisioner struct {
}

// Deprovision is the Deprovisioner interface implementation. It returns nil, nil
func (d *Deprovisioner) Deprovision(instanceID, serviceID, planID string) (*mode.DeprovisionResponse, error) {
	return nil, nil
}
