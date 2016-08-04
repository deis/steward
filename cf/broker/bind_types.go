package broker

// JSON types for requests to the proxy and backend CF broker

type backendBindingRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

type namespaceParams struct {
	TargetNamespace string `json:"target_namespace"`
}

type fullBindingRequest struct {
	BackendReq backendBindingRequest `json:",inline"`
	Parameters namespaceParams       `json:"parameters"`
}

// JSON types for responses from the backend CF broker and the proxy

// the credentials sent in the CF broker response
type bindingCredentials struct {
	URI      string `json:"uri"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type backendBindingResponse struct {
	Credentials bindingCredentials `json:"credentials"`
}

type fullBindingResponse struct {
	ConfigMapAndSecret configMapAndSecret `json:"credentials"`
}

// a pair of qualified names describing the config map and associated secrets
type configMapAndSecret struct {
	ConfigMapInfo *qualifiedName   `json:"config_map_info"`
	SecretsInfo   []*qualifiedName `json:"secrets_info"`
}

// a (name, namespace) pair to identify exactly where a resource is
type qualifiedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// // the entire request sent to the proxy
// type bindingRequest struct {
// 	ServiceID  string           `json:"service_id"`
// 	PlanID     string           `json:"plan_id"`
// 	Parameters *namespaceParams `json:"parameters"`
// }
//
// type backendBindingResponse struct {
// 	Creds bindingCredentials `json:"credentials"`
// }
//
// type proxyRespBindingCredentials struct {
// 	Creds              bindingCredentials `json:",inline"`
// 	ConfigMapAndSecret configMapAndSecret `json:",inline"`
// }
//
// // the entire response sent from the proxy
// type fullBindingResponse struct {
// 	Creds proxyRespBindingCredentials `json:"credentials"`
// }
