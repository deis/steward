package k8s

const (
	// ActionProvision is the action indicating that a service should be provision
	ActionProvision ServicePlanClaimAction = "provision"
	// ActionBind is the action indicating that a service is should be bound
	ActionBind ServicePlanClaimAction = "bind"
	// ActionUnbind is the action indicating that a servie should be unbound
	ActionUnbind ServicePlanClaimAction = "unbind"
	// ActionDeprovision is the action indicating that a service should be deprovisioned
	ActionDeprovision ServicePlanClaimAction = "deprovision"
	// ActionCreate is the action indicating that a service should be provisioned and bound, in that order, in the same operation
	ActionCreate ServicePlanClaimAction = "create"
	// ActionDelete is the actions indicating that a service should be unbound and deprovisioned, in that order in the same operation
	ActionDelete ServicePlanClaimAction = "delete"
)

// ServicePlanClaimAction is the type representing the current action a consumer has requested on a claim. It implements fmt.Stringer
type ServicePlanClaimAction string

// StringIsServicePlanClaimAction returns true if s == a.String()
func StringIsServicePlanClaimAction(s string, a ServicePlanClaimAction) bool {
	return s == a.String()
}

// String is the fmt.Stringer implementation
func (a ServicePlanClaimAction) String() string {
	return string(a)
}
