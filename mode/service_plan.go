package mode

// ServicePlan is the steward representation of a service plan. It's also compatible with the CloudFoundtry service plan object. See https://docs.cloudfoundry.org/services/api.html#PObject for more detail
type ServicePlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Metadata ServiceMetadata `json:"metadata"`
	Free bool `json:"free"`
}
