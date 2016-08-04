package broker

import (
	"fmt"
	"net/http"

	"github.com/deis/steward/cf"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

const (
	statusUnprocessableEntity = 422
	instanceIDPathKey         = "instance_id"
	bindingIDPathKey          = "binding_id"
)

// Handler returns the HTTP handler for all CloudFoundry API endpoints
func Handler(
	logger loggo.Logger,
	cl *cf.Client,
	frontendAuth,
	backendAuth *web.BasicAuth,
	cmCreator k8s.ConfigMapCreator,
	secCreator k8s.SecretCreator,
) http.Handler {

	r := mux.NewRouter()
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		provisioningHandler(logger, cl, frontendAuth, backendAuth),
	).Methods("PUT")
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}/service_bindings/{%s}", instanceIDPathKey, bindingIDPathKey),
		bindingHandler(logger, cl, frontendAuth, backendAuth, cmCreator, secCreator),
	).Methods("PUT")
	return r
}
