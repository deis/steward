package fake

import (
	"github.com/deis/steward/mode"
)

type DeprovisionCall struct {
	InstanceID string
	ServiceID  string
	PlanID     string
}

// Deprovisioner is a fake implementation of (github.com/deis/steward/mode).Deprovisioner, suitable for usage in unit tests
type Deprovisioner struct {
	Deprovisions []DeprovisionCall
}

// Deprovision is the Deprovisioner interface implementation. It packages the function parameters into a DeprovisionCall, appends it to d.Deprovisons, and returns nil, nil. This function is not concurrency safe
func (d *Deprovisioner) Deprovision(instanceID string, dReq *mode.DeprovisionRequest) (*mode.DeprovisionResponse, error) {
	d.Deprovisions = append(d.Deprovisions, DeprovisionCall{
		InstanceID: instanceID,
		ServiceID:  dReq.ServiceID,
		PlanID:     dReq.PlanID,
	})
	return nil, nil
}
