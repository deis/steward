package fake

import (
	"github.com/deis/steward/mode"
)

// Binder is a fake (github.com/deis/steward/mode).Binder implementation, suitable for use in unit tests
type Binder struct {
}

// Bind is the Binder interface implementation. It returns nil, nil
func (b *Binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	return nil, nil
}
