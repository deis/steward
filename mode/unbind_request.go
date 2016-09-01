package mode

// UnbindRequest represents a request to unbind from a service. It is marked with JSON struct tags so that it can be encoded to, and decoded from the CloudFoundry unbinding request body format. See https://docs.cloudfoundry.org/services/api.html#unbinding for more details
type UnbindRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}
