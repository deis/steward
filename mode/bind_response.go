package mode

// BindResponse is the response to a binding request
type BindResponse struct {
	Status int
	Creds  JSONObject
}
