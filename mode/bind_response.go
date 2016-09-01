package mode

// BindResponse represents a response to a binding request. It is marked with JSON struct tags so that it can be encoded to, and decoded from the CloudFoundry binding response body format. See https://docs.cloudfoundry.org/services/api.html#binding for more details
type BindResponse struct {
	Creds JSONObject `json:"credentials"`
}
