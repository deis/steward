package mode

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
