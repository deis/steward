package brokerapi

import (
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	serviceIDQueryKey  = "service_id"
	planIDQueryKey     = "plan_id"
	instanceIDPathKey  = "instance_id"
	bindingIDPathKey   = "binding_id"
	targetNamespaceKey = "target_namespace"
)

// Handler returns the HTTP handler for all CloudFoundry API endpoints
func Handler(
	cataloger mode.Cataloger,
	lifecycler mode.Lifecycler,
	frontendAuth *web.BasicAuth,
	cmNamespacer kcl.ConfigMapsNamespacer,
) http.Handler {

	r := mux.NewRouter()

	// provisioning
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		web.WithBasicAuth(
			frontendAuth,
			provisioningHandler(lifecycler),
		),
	).Methods("PUT")

	// deprovisioning
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		web.WithBasicAuth(
			frontendAuth,
			deprovisioningHandler(lifecycler),
		),
	).Methods("DELETE")

	// binding
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}/service_bindings/{%s}", instanceIDPathKey, bindingIDPathKey),
		web.WithBasicAuth(
			frontendAuth,
			bindingHandler(lifecycler, cmNamespacer),
		),
	).Methods("PUT")

	// unbinding
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}/service_bindings/{%s}", instanceIDPathKey, bindingIDPathKey),
		web.WithBasicAuth(
			frontendAuth,
			unbindHandler(lifecycler, cmNamespacer),
		),
	)

	// catalog listing
	r.Handle(
		"/v2/catalog",
		web.WithBasicAuth(frontendAuth,
			catalogHandler(cataloger),
		),
	).Methods("GET")

	return r
}
