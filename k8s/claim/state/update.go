package state

import (
	"fmt"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
)

// Update represents the update to the state of a ServicePlanClaim
type Update interface {
	fmt.Stringer
	Status() k8s.ServicePlanClaimStatus
	Description() string
	InstanceID() string
	BindID() string
	Extra() mode.JSONObject
}

// UpdateClaim updates claim in-place, according to update
func UpdateClaim(claim *k8s.ServicePlanClaim, update Update) {
	switch u := update.(type) {
	case statusUpdate:
		claim.Status = u.status.String()
	case errUpdate:
		claim.Status = k8s.StatusFailed.String()
		claim.StatusDescription = u.err.Error()
	case fullUpdate:
		claim.Status = u.status.String()
		claim.StatusDescription = u.description
		if len(u.instanceID) > 0 {
			claim.InstanceID = u.instanceID
		}
		if len(u.bindID) > 0 {
			claim.BindID = u.bindID
		}
		if len(u.extra) > 0 {
			claim.Extra = u.extra
		}
	default:
	}
}

// UpdateIsTerminal returns true if u will, after applied to a ServicePlanClaim, result in the claim being in a potentially terminal state. Note that "potentially terminal state" doesn't necessarily mean that the claim is no longer actionable. It just means that steward has the option to not _automatically_ take more action on the claim. For example, a 'provision' claim will result in the claim becoming 'provisioned', which is a potentially terminal state. At this point, steward doesn't need to take further action on the claim, but it will if the user sets the claim's action to 'bind'
func UpdateIsTerminal(u Update) bool {
	switch u.Status() {
	case k8s.StatusFailed, k8s.StatusBound, k8s.StatusProvisioned, k8s.StatusUnbound, k8s.StatusDeprovisioned:
		return true
	default:
		return false
	}
}
