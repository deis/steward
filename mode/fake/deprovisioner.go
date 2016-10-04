package fake

import (
	"github.com/deis/steward/mode"
)

type DeprovisionCall struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	Params     mode.JSONObject
}

// Deprovisioner is a fake implementation of (github.com/deis/steward/mode).Deprovisioner, suitable for usage in unit tests
type Deprovisioner struct {
	Deprovisions []DeprovisionCall
	Resp         *mode.DeprovisionResponse
	Err          error
}

// Deprovision is the Deprovisioner interface implementation. It packages the function parameters into a DeprovisionCall, appends it to d.Deprovisons, and returns d.Resp, d.Err. This function is not concurrency safe
func (d *Deprovisioner) Deprovision(instanceID string, dReq *mode.DeprovisionRequest) (*mode.DeprovisionResponse, error) {
	d.Deprovisions = append(d.Deprovisions, DeprovisionCall{
		InstanceID: instanceID,
		ServiceID:  dReq.ServiceID,
		PlanID:     dReq.PlanID,
		Params:     dReq.Parameters,
	})
	return d.Resp, d.Err
}
