package state

import (
	"github.com/deis/steward/k8s"
	"k8s.io/client-go/1.4/pkg/watch"
)

// Current represents the current state of a ServicePlanClaim
type Current struct {
	Status      k8s.ServicePlanClaimStatus
	StatusValid bool
	Action      k8s.ServicePlanClaimAction
	EventType   watch.EventType
}

// NewCurrent returns a new Current with the given parameters and StatusValid set to true
func NewCurrent(
	status k8s.ServicePlanClaimStatus,
	action k8s.ServicePlanClaimAction,
	evtType watch.EventType,
) Current {
	return Current{
		Status:      status,
		StatusValid: true,
		Action:      action,
		EventType:   evtType,
	}
}

// NewCurrentNoStatus returns a new Current without a Status and StatusValid set to false
func NewCurrentNoStatus(action k8s.ServicePlanClaimAction, evtType watch.EventType) Current {
	return Current{
		Action:      action,
		StatusValid: false,
		EventType:   evtType,
	}
}
