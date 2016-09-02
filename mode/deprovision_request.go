package mode

// DeprovisionRequest represents a request to do a service deprovision operation. This struct is JSON-compatible with the request body detailed at https://docs.cloudfoundry.org/services/api.html#deprovisioning
type DeprovisionRequest struct {
	ServiceID  string     `json:"service_id"`
	PlanID     string     `json:"plan_id"`
	Parameters JSONObject `json:parameters,omitempty`
}
