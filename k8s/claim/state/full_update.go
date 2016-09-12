package state

import (
	"fmt"

	"github.com/deis/steward/mode"
)

type fullUpdate struct {
	status      mode.Status
	description string
	instanceID  string
	bindID      string
	extra       mode.JSONObject
}

// FullUpdate returns an Update implementation with all fields filled in
func FullUpdate(st mode.Status, desc, instID, bindID string, extra mode.JSONObject) Update {
	return fullUpdate{
		status:      st,
		description: desc,
		instanceID:  instID,
		bindID:      bindID,
		extra:       extra,
	}
}

func (f fullUpdate) String() string {
	return fmt.Sprintf(
		"full update. status = %s, description = '%s', instanceID = %s, bindID = %s, extra = %s",
		f.status,
		f.description,
		f.instanceID,
		f.bindID,
		f.extra,
	)
}

func (f fullUpdate) Status() mode.Status {
	return f.status
}

func (f fullUpdate) Description() string {
	return f.description
}

func (f fullUpdate) InstanceID() string {
	return f.instanceID
}
func (f fullUpdate) BindID() string {
	return f.bindID
}
func (f fullUpdate) Extra() mode.JSONObject {
	return f.extra
}
