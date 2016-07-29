package cf

// Service is the represntation of a CloudFoundry service. See https://docs.cloudfoundry.org/services/api.html#catalog-mgmt for more detail
type Service struct {
	ServiceInfo
	Plans []ServicePlan `json:"plans"`
}

// ServiceInfo represents all of the information about a service except for its plans
type ServiceInfo struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	// Tags          []string         `json:"tags"`
	// Requires      []string         `json:"requires"`
	// Bindable      bool             `json:"bindable"`
	// Metadata      ServicesMetadata `json:"metadata"`
	PlanUpdatable bool `json:"plan_updateable"`
}

// ServiceMetadata is the representation of a CloudFoundry service metadata. See https://docs.cloudfoundry.org/services/catalog-metadata.html for more detail
type ServiceMetadata struct {
}
