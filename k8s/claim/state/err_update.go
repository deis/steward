package state

import (
	"fmt"

	"github.com/deis/steward/mode"
)

// ErrUpdate is an Update implementation that sets the claim to a failed state
type errUpdate struct {
	err error
}

// ErrUpdate returns a new Update implementation that has a failed status and status description equal to e.Error()
func ErrUpdate(e error) Update {
	return errUpdate{err: e}
}

func (e errUpdate) String() string {
	return fmt.Sprintf("status update to failure with error %s", e.err)
}

func (e errUpdate) Status() mode.Status {
	return mode.StatusFailed
}

func (e errUpdate) Description() string {
	return e.err.Error()
}

func (e errUpdate) InstanceID() string {
	return ""
}
func (e errUpdate) BindID() string {
	return ""
}
func (e errUpdate) Extra() mode.JSONObject {
	return mode.EmptyJSONObject()
}
