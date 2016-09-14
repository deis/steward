package helm

import (
	"github.com/deis/steward/mode"
)

type unbinder struct{}

func (u unbinder) Unbind(instanceID, bindingID string, unbindRequest *mode.UnbindRequest) error {
	return nil
}

// newUnbinder returns a Tiller-backed mode.Unbinder
func newUnbinder() mode.Unbinder {
	return unbinder{}
}
