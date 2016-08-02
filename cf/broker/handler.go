package broker

import (
	"fmt"
	"net/http"

	"github.com/deis/steward/cf"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

const (
	statusUnprocessableEntity = 422
)

// Handler returns the HTTP handler for all CloudFoundry API endpoints
func Handler(logger loggo.Logger, cl *cf.Client, frontendAuth, backendAuth *web.BasicAuth) http.Handler {
	r := mux.NewRouter()
	r.Handle(
		fmt.Sprintf("/v2/service_instances/{%s}", instanceIDPathKey),
		provisioningHandler(logger, cl, frontendAuth, backendAuth),
	).Methods("PUT")
	return r
}
