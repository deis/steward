package mode

import (
	"fmt"
)

// Action is the type representing the current action a consumer has requested on a claim. It implements fmt.Stringer
type Action string

// StringIsAction returns true if s == a.String()
func StringIsAction(s string, a Action) bool {
	return s == a.String()
}

// String is the fmt.Stringer implementation
func (a Action) String() string {
	return string(a)
}

// Status is the type representing the current status of a claim. It implements fmt.Stringer
type Status string

// StringIsStatus returns true if s == st.String()
func StringIsStatus(s string, st Status) bool {
	return s == st.String()
}

// String is the fmt.Stringer interface implementation
func (s Status) String() string {
	return string(s)
}

const (
	serviceIDKey            = "service-id"
	planIDKey               = "plan-id"
	claimIDMapKey           = "claim-id"
	actionMapKey            = "action"
	statusMapKey            = "status"
	statusDescriptionMapKey = "status-description"
	targetNameMapKey        = "target-name"
	targetNamespaceMapKey   = "target-namespace"
	instanceIDMapKey        = "instance-id"
	bindIDMapKey            = "bind-id"

	// ActionProvision is the action indicating that a service should be provision
	ActionProvision Action = "provision"
	// ActionBind is the action indicating that a service is should be bound
	ActionBind Action = "bind"
	// ActionUnbind is the action indicating that a servie should be unbound
	ActionUnbind Action = "unbind"
	// ActionDeprovision is the action indicating that a service should be deprovisioned
	ActionDeprovision Action = "deprovision"

	// StatusProvisioning is the status indicating that the provisioning process has started
	StatusProvisioning Status = "provisioning"
	// StatusProvisioned is the status indicating that the provisioning process has succeeded
	StatusProvisioned Status = "provisioned"
	// StatusBinding is the status indicating that the binding process has started
	StatusBinding Status = "binding"
	// StatusBound is the status indicating that the binding process has succeeded
	StatusBound Status = "bound"
	// StatusUnbinding is the status indicating that the unbinding process has started
	StatusUnbinding Status = "unbinding"
	// StatusUnbound is the status indicating that the unbinding process has succeeded
	StatusUnbound Status = "unbound"
	// StatusDeprovisioning is the status indicating that the deprovisioning process has started
	StatusDeprovisioning Status = "deprovisioning"
	// StatusDeprovisioned is the status indicating that the deprovisioning process has succeeded
	StatusDeprovisioned Status = "deprovisioned"
	// StatusFailed is the status indicating that a service's creation or deletion operation has failed for some reason. The human-readable explanation of the failure will be written to the status description
	StatusFailed Status = "failed"
)

type errServicePlanClaimMapMissingKey struct {
	key string
}

func (e errServicePlanClaimMapMissingKey) Error() string {
	return fmt.Sprintf("map to convert to service plan claim is missing key %s", e.key)
}

// ServicePlanClaim is the json-encodable struct that represents a service plan claim. See https://github.com/deis/steward/blob/master/DATA_STRUCTURES.md#serviceplanclaim for more detail. This struct implements fmt.Stringer
type ServicePlanClaim struct {
	TargetName        string `json:"target-name"`
	TargetNamespace   string `json:"target-namespace"`
	ServiceID         string `json:"service-id"`
	PlanID            string `json:"plan-id"`
	ClaimID           string `json:"claim-id"`
	Action            string `json:"action"`
	Status            string `json:"status"`
	StatusDescription string `json:"status-description"`
	InstanceID        string `json:"instance-id"`
	BindID            string `json:"bind-id"`
}

// ServicePlanClaimFromMap attempts to convert m to a ServicePlanClaim. If the map was malformed or missing any keys, returns nil and an appropriate error
func ServicePlanClaimFromMap(m map[string]string) (*ServicePlanClaim, error) {
	targetName, ok := m[targetNameMapKey]
	if !ok {
		return nil, errServicePlanClaimMapMissingKey{key: targetNameMapKey}
	}
	targetNamespace, ok := m[targetNamespaceMapKey]
	if !ok {
		return nil, errServicePlanClaimMapMissingKey{key: targetNamespaceMapKey}
	}
	serviceID, ok := m[serviceIDKey]
	if !ok {
		return nil, errServicePlanClaimMapMissingKey{key: serviceIDKey}
	}
	planID, ok := m[planIDKey]
	if !ok {
		return nil, errServicePlanClaimMapMissingKey{key: planIDKey}
	}
	claimID, ok := m[claimIDMapKey]
	if !ok {
		return nil, errServicePlanClaimMapMissingKey{key: claimIDMapKey}
	}
	action, ok := m[actionMapKey]
	if !ok {
		return nil, errServicePlanClaimMapMissingKey{key: actionMapKey}
	}
	// the following fields may be empty when the application submits them, so don't error if they're missing
	status := m[statusMapKey]
	statusDescription := m[statusDescriptionMapKey]
	instanceID := m[instanceIDMapKey]
	bindID := m[bindIDMapKey]

	return &ServicePlanClaim{
		TargetName:        targetName,
		TargetNamespace:   targetNamespace,
		ServiceID:         serviceID,
		PlanID:            planID,
		ClaimID:           claimID,
		Action:            action,
		Status:            status,
		StatusDescription: statusDescription,
		InstanceID:        instanceID,
		BindID:            bindID,
	}, nil
}

// ToMap returns s represented as a map[string]strinrg
func (s ServicePlanClaim) ToMap() map[string]string {
	return map[string]string{
		targetNameMapKey:        s.TargetName,
		targetNamespaceMapKey:   s.TargetNamespace,
		serviceIDKey:            s.ServiceID,
		planIDKey:               s.PlanID,
		claimIDMapKey:           s.ClaimID,
		actionMapKey:            s.Action,
		statusMapKey:            s.Status,
		statusDescriptionMapKey: s.StatusDescription,
		instanceIDMapKey:        s.InstanceID,
		bindIDMapKey:            s.BindID,
	}
}

// String is the fmt.Stringer interface implementation
func (s ServicePlanClaim) String() string {
	return fmt.Sprintf("%s", s.ToMap())
}
