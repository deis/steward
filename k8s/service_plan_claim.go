package k8s

import (
	"fmt"

	"github.com/deis/steward/mode"
)

const (
	serviceIDKey            = "service-id"
	planIDKey               = "plan-id"
	claimIDMapKey           = "claim-id"
	actionMapKey            = "action"
	statusMapKey            = "status"
	statusDescriptionMapKey = "status-description"
	targetNameMapKey        = "target-name"
	instanceIDMapKey        = "instance-id"
	bindIDMapKey            = "bind-id"
	extraMapKey             = "extra"
)

type errDataMapMissingKey struct {
	key string
}

func (e errDataMapMissingKey) Error() string {
	return fmt.Sprintf("map to convert to service plan claim is missing key %s", e.key)
}

// ServicePlanClaim is the json-encodable struct that represents a service plan claim. See https://github.com/deis/steward/blob/master/DATA_STRUCTURES.md#serviceplanclaim for more detail. This struct implements fmt.Stringer
type ServicePlanClaim struct {
	TargetName        string          `json:"target-name"`
	ServiceID         string          `json:"service-id"`
	PlanID            string          `json:"plan-id"`
	ClaimID           string          `json:"claim-id"`
	Action            string          `json:"action"`
	Status            string          `json:"status"`
	StatusDescription string          `json:"status-description"`
	InstanceID        string          `json:"instance-id"`
	BindID            string          `json:"bind-id"`
	Extra             mode.JSONObject `json:"extra"`
}

// ServicePlanClaimFromMap attempts to convert m to a ServicePlanClaim. If the map was malformed or missing any keys, returns nil and an appropriate error
func ServicePlanClaimFromMap(m map[string]string) (*ServicePlanClaim, error) {
	targetName, ok := m[targetNameMapKey]
	if !ok {
		return nil, errDataMapMissingKey{key: targetNameMapKey}
	}
	serviceID, ok := m[serviceIDKey]
	if !ok {
		return nil, errDataMapMissingKey{key: serviceIDKey}
	}
	planID, ok := m[planIDKey]
	if !ok {
		return nil, errDataMapMissingKey{key: planIDKey}
	}
	claimID, ok := m[claimIDMapKey]
	if !ok {
		return nil, errDataMapMissingKey{key: claimIDMapKey}
	}
	action, ok := m[actionMapKey]
	if !ok {
		return nil, errDataMapMissingKey{key: actionMapKey}
	}
	// the following fields may be empty when the application submits them, so don't error if they're missing
	status := m[statusMapKey]
	statusDescription := m[statusDescriptionMapKey]
	instanceID := m[instanceIDMapKey]
	bindID := m[bindIDMapKey]
	extraStr := m[extraMapKey]
	extra, err := mode.JSONObjectFromString(extraStr)
	if err != nil {
		return nil, err
	}

	return &ServicePlanClaim{
		TargetName:        targetName,
		ServiceID:         serviceID,
		PlanID:            planID,
		ClaimID:           claimID,
		Action:            action,
		Status:            status,
		StatusDescription: statusDescription,
		InstanceID:        instanceID,
		BindID:            bindID,
		Extra:             extra,
	}, nil
}

// ToMap returns s represented as a map[string]strinrg
func (s ServicePlanClaim) ToMap() map[string]string {
	return map[string]string{
		targetNameMapKey:        s.TargetName,
		serviceIDKey:            s.ServiceID,
		planIDKey:               s.PlanID,
		claimIDMapKey:           s.ClaimID,
		actionMapKey:            s.Action,
		statusMapKey:            s.Status,
		statusDescriptionMapKey: s.StatusDescription,
		instanceIDMapKey:        s.InstanceID,
		bindIDMapKey:            s.BindID,
		extraMapKey:             s.Extra.EncodeToString(),
	}
}

// String is the fmt.Stringer interface implementation
func (s ServicePlanClaim) String() string {
	return fmt.Sprintf("%s", s.ToMap())
}
