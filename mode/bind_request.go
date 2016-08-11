package mode

const (
	// TargetNamespaceKey is the required key for the target namespace value
	TargetNamespaceKey = "target_namespace"
	// TargetNameKey is the required key for the target name value
	TargetNameKey = "target_name"
)

// BindRequest represents a request to bind to a service. It is marked with JSON struct tags so that it can be encoded to, and decoded from the CloudFoundry binding request body format. See https://docs.cloudfoundry.org/services/api.html#binding for more details
type BindRequest struct {
	ServiceID  string     `json:"service_id"`
	PlanID     string     `json:"plan_id"`
	Parameters JSONObject `json:"parameters"`
}

// TargetNamespace returns the target namespace in b.Parameters, or an error if it's missing
func (b BindRequest) TargetNamespace() (string, error) {
	return b.Parameters.String(TargetNamespaceKey)
}

// TargetName returns the target name in b.Parameters, or an error if it's missing
func (b BindRequest) TargetName() (string, error) {
	return b.Parameters.String(TargetNameKey)
}
