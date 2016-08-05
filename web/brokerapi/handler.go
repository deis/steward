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
	instanceIDPathKey  = "instance_id"
	bindingIDPathKey   = "binding_id"
	targetNamespaceKey = "target_namespace"
)

// Handler returns the HTTP handler for all CloudFoundry API endpoints
func Handler(
	logger loggo.Logger,
	cataloger mode.Cataloger,
	provisioner mode.Provisioner,
	binder mode.Binder,
	frontendAuth *web.BasicAuth,
	cmCreator k8s.ConfigMapCreator,
	secCreator k8s.SecretCreator,
) http.Handler {

	r := mux.NewRouter()
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		provisioningHandler(logger, provisioner, frontendAuth),
	).Methods("PUT")
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}/service_bindings/{%s}", instanceIDPathKey, bindingIDPathKey),
		bindingHandler(logger, binder, frontendAuth, cmCreator, secCreator),
	).Methods("PUT")
	r.Handle("/v2/catalog", catalogHandler(logger, cataloger, frontendAuth)).Methods("GET")
	return r
}
