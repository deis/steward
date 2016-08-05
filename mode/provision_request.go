package mode

// ProvisionRequest represents a request to do a service provision operation. This struct is JSON-compatible with the request body detailed at https://docs.cloudfoundry.org/services/api.html#provisioning
type ProvisionRequest struct {
	OrganizationGUID string     `json:"organization_guid"`
	PlanID           string     `json:"plan_id"`
	ServiceID        string     `json:"service_id"`
	SpaceGUID        string     `json:"space_guid"`
	Parameters       JSONObject `json:"parameters"`
}
