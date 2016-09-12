package state

import (
	"fmt"

	"github.com/deis/steward/mode"
)

type statusUpdate struct {
	status mode.Status
}

// StatusUpdate returns an Update implementation that only updates the status field of a claim
func StatusUpdate(st mode.Status) Update {
	return statusUpdate{status: st}
}

func (s statusUpdate) String() string {
	return fmt.Sprintf("status update to %s", s.Status)
}

func (s statusUpdate) Status() mode.Status {
	return s.status
}

func (s statusUpdate) Description() string {
	return ""
}

func (s statusUpdate) InstanceID() string {
	return ""
}

func (s statusUpdate) BindID() string {
	return ""
}

func (s statusUpdate) Extra() mode.JSONObject {
	return mode.EmptyJSONObject()
}
