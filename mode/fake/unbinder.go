package fake

import (
	"github.com/deis/steward/mode"
)

type UnbindCall struct {
	InstanceID    string
	BindID        string
	UnbindRequest *mode.UnbindRequest
}

// Unbinder is a fake (github.com/deis/steward/mode).Unbinder implementation, suitable for use in unit tests
type Unbinder struct {
	UnbindCalls []*UnbindCall
}

// Unbind is the Unbinder interface implementaion. It returns nil
func (u *Unbinder) Unbind(instanceID, bindingID string, uReq *mode.UnbindRequest) error {
	u.UnbindCalls = append(u.UnbindCalls, &UnbindCall{
		InstanceID:    instanceID,
		BindID:        bindingID,
		UnbindRequest: uReq,
	})
	return nil
}
