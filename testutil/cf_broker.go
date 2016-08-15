package testutil

import (
	"encoding/json"
	"net/http"

	"github.com/arschles/testsrv"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
)

// NewCFBroker returns a new Server that can handle requests to
func NewCFBroker(creds *web.BasicAuth, bindCreds map[string]string) *testsrv.Server {
	r := mux.NewRouter()
	r.Handle(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
		web.WithBasicAuth(creds, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			retMap := map[string]interface{}{
				"credentials": bindCreds,
			}
			if err := json.NewEncoder(w).Encode(retMap); err != nil {
				http.Error(w, "error encoding JSON", http.StatusInternalServerError)
				return
			}
		})),
	)
	return testsrv.StartServer(r)
}
