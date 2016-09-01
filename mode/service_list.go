package mode

// ServiceList is a wrapper for a list of services and represents a response from a request to list services provided by a broker. It is marked with JSON struct tags so that it can be encoded to, and decoded from the CloudFoundry catalog list response body format. See https://docs.cloudfoundry.org/services/api.html#catalog-mgmt
type ServiceList struct {
	Services []*Service `json:"services"`
}
