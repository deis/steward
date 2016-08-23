package fake

import (
	"github.com/deis/steward/mode"
)

type BindCall struct {
	InstanceID string
	BindingID  string
	Req        *mode.BindRequest
}

// Binder is a fake (github.com/deis/steward/mode).Binder implementation, suitable for use in unit tests
type Binder struct {
	Binds []BindCall
	Res   *mode.BindResponse
}

// Bind is the Binder interface implementation. It constructs a new BindCall from the function params, then returns b.Res, nil. This function is not concurrency safe
func (b *Binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	b.Binds = append(b.Binds, BindCall{
		InstanceID: instanceID,
		BindingID:  bindingID,
		Req:        bindRequest,
	})
	return b.Res, nil
}
