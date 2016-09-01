package state

import (
	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/watch"
)

// Current represents the current state of a ServicePlanClaim
type Current struct {
	Status      mode.Status
	StatusValid bool
	Action      mode.Action
	EventType   watch.EventType
}

// NewCurrent returns a new Current with the given parameters and StatusValid set to true
func NewCurrent(status mode.Status, action mode.Action, evtType watch.EventType) Current {
	return Current{
		Status:      status,
		StatusValid: true,
		Action:      action,
		EventType:   evtType,
	}
}

// NewCurrentNoStatus returns a new Current without a Status and StatusValid set to false
func NewCurrentNoStatus(action mode.Action, evtType watch.EventType) Current {
	return Current{
		Action:      action,
		StatusValid: false,
		EventType:   evtType,
	}
}
