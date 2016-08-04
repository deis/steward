package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deis/steward/cf"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func provisioningHandler(logger loggo.Logger, cl *cf.Client, frontendAuth, backendAuth *web.BasicAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceID, ok := vars[instanceIDPathKey]
		if !ok {
			http.Error(w, "missing instance ID", http.StatusBadRequest)
			return
		}
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authorization missing", http.StatusBadRequest)
			return
		}
		if username != frontendAuth.Username || password != frontendAuth.Password {
			http.Error(w, "wrong login credentials", http.StatusUnauthorized)
			return
		}
		req, err := cl.Put(logger, r.Body, "/v2", "service_instances", instanceID)
		if err != nil {
			logger.Debugf("error creating PUT request to backend CF broker (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := cl.Client.Do(req)
		if err != nil {
			logger.Debugf("error making request to backend service broker (%s)", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res.StatusCode == http.StatusConflict || res.StatusCode == statusUnprocessableEntity {
			logger.Debugf("response code from backend service broker was %d", res.StatusCode)
			http.Error(w, fmt.Sprintf("error code from backend service broker %d", res.StatusCode), http.StatusInternalServerError)
			return
		}
		resp := new(backendProvisionResp)
		if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
			logger.Debugf("error decoding JSON from backend service broker (%s)", err)
			http.Error(w, "error decoding response from backend service broker", http.StatusInternalServerError)
			return
		}
		respStr := fmt.Sprintf(`{"operation":"%s"}`, resp.Operation)
		w.WriteHeader(res.StatusCode)
		w.Write([]byte(respStr))
	})
}

type backendProvisionResp struct {
	Operation string `json:"operation"`
}
