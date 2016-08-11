package brokerapi

import (
	"fmt"
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
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
	logger loggo.Logger,
	cataloger mode.Cataloger,
	provisioner mode.Provisioner,
	deprovisioner mode.Deprovisioner,
	binder mode.Binder,
	unbinder mode.Unbinder,
	frontendAuth *web.BasicAuth,
	cmCreatorDeleter k8s.ConfigMapCreatorDeleter,
) http.Handler {

	r := mux.NewRouter()

	// provisioning
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		withBasicAuth(
			frontendAuth,
			provisioningHandler(logger, provisioner),
		),
	).Methods("PUT")

	// deprovisioning
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		withBasicAuth(
			frontendAuth,
			deprovisioningHandler(logger, deprovisioner),
		),
	).Methods("DELETE")

	// binding
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}/service_bindings/{%s}", instanceIDPathKey, bindingIDPathKey),
		withBasicAuth(
			frontendAuth,
			bindingHandler(logger, binder, cmCreatorDeleter),
		),
	).Methods("PUT")

	// unbinding
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}/service_bindings/{%s}", instanceIDPathKey, bindingIDPathKey),
		withBasicAuth(
			frontendAuth,
			unbindHandler(logger, unbinder, cmCreatorDeleter),
		),
	)

	// catalog listing
	r.Handle(
		"/v2/catalog",
		withBasicAuth(frontendAuth,
			catalogHandler(logger, cataloger),
		),
	).Methods("GET")

	return r
}
