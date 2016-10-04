package mode

// ProvisionResponse represents a response to a provisioning request. It is marked with JSON struct tags so that it can be encoded to, and decoded from the CloudFoundry provisioning response body format. See https://docs.cloudfoundry.org/services/api.html#provisioning for more details
type ProvisionResponse struct {
	Operation string     `json:"operation"`
	IsAsync   bool       `json:"-"`
	Extra     JSONObject `json:"extra,omitempty"`
}
