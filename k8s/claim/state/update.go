package state

import (
	"fmt"

	"github.com/deis/steward/mode"
)

// Update represents the update to the state of a ServicePlanClaim
type Update struct {
	NewStatus            mode.Status
	NewStatusDescription string
	NewExtra             mode.JSONObject
}

// NewUpdate creates a new Update from the given parameters
func NewUpdate(stat mode.Status, descr string, extra mode.JSONObject) Update {
	return Update{
		NewStatus:            stat,
		NewStatusDescription: descr,
		NewExtra:             extra,
	}
}

// UpdateClaim updates claim in-place, according to update
func UpdateClaim(claim *mode.ServicePlanClaim, update Update) {
	claim.Status = update.NewStatus.String()
	claim.StatusDescription = update.NewStatusDescription
	claim.Extra = update.NewExtra
}

// ErrUpdate creates a new Update with new status failed, status description the value of err.Error(), and the new extra field set to extra
func ErrUpdate(err error, extra mode.JSONObject) Update {
	return Update{
		NewStatus:            mode.StatusFailed,
		NewStatusDescription: err.Error(),
		NewExtra:             extra,
	}
}

func (u Update) String() string {
	return fmt.Sprintf("status update status = %s, descr = '%s', extra = '%+v'", u.NewStatus, u.NewStatusDescription, u.NewExtra)
}

// IsTerminal returns true if u will, after applied to a ServicePlanClaim, result in the claim being in a potentially terminal state. Note that "potentially terminal state" doesn't necessarily mean that the claim is no longer actionable. It just means that steward has the option to not _automatically_ take more action on the claim. For example, a 'provision' claim will result in the claim becoming 'provisioned', which is a potentially terminal state. At this point, steward doesn't need to take further action on the claim, but it will if the user sets the claim's action to 'bind'
func (u Update) IsTerminal() bool {
	switch u.NewStatus {
	case mode.StatusFailed, mode.StatusBound, mode.StatusProvisioned, mode.StatusUnbound, mode.StatusDeprovisioned:
		return true
	default:
		return false
	}
}
