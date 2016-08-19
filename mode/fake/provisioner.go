package fake

import (
	"github.com/deis/steward/mode"
)

// Provisioner is a fake implementation of (github.com/deis/steward/mode).Provisioner, suitable for usage in unit tests
type Provisioner struct {
}

// Provision is the Provisioner interface implementation. It returns nil, nil
func (p *Provisioner) Provision(instanceID string, req *mode.ProvisionRequest) (*mode.ProvisionResponse, error) {
	return nil, nil
}
