package broker

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/deis/steward/cf"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func bindingHandler(
	logger loggo.Logger,
	cl *cf.Client,
	frontendAuth,
	backendAuth *web.BasicAuth,
	cmCreator k8s.ConfigMapCreator,
	secCreator k8s.SecretCreator,
) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceID, ok := vars[instanceIDPathKey]
		if !ok {
			http.Error(w, "missing instance ID", http.StatusBadRequest)
			return
		}
		bindingID, ok := vars[bindingIDPathKey]
		if !ok {
			http.Error(w, "missing binding ID", http.StatusBadRequest)
			return
		}
		fullReq := new(fullBindingRequest)
		if err := json.NewDecoder(r.Body).Decode(fullReq); err != nil {
			logger.Debugf("error decoding the full broker request (%s)", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		rdr := new(bytes.Buffer)
		if err := json.NewEncoder(rdr).Encode(fullReq.BackendReq); err != nil {
			logger.Debugf("error encoding the backend broker request (%s)", err)
			http.Error(w, "error encoding the backend broker request", http.StatusInternalServerError)
			return
		}

		res, err := cl.DoPut(logger, rdr, "v2", "service_instances", instanceID, "service_bindings", bindingID)
		if err != nil {
			logger.Debugf("error executing PUT request to backend broker (%s)", err)
			http.Error(w, "error executing backend request", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusConflict || res.StatusCode == statusUnprocessableEntity {
			logger.Debugf("backend CF broker request returned error code %d", res.StatusCode)
			http.Error(w, "backend CF broker returned error code", http.StatusInternalServerError)
			return
		}
		resp := new(backendBindingResponse)
		if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
			logger.Debugf("backend CF broker returned malformed response body (%s)", err)
			http.Error(w, "backend CF broker returned malformed response body", http.StatusInternalServerError)
			return
		}

		configMapQualified, secretsQualified, err := writeToKubernetes(
			fullReq.BackendReq.ServiceID,
			fullReq.BackendReq.PlanID,
			fullReq.Parameters.TargetNamespace,
			resp.Credentials,
			cmCreator,
			secCreator,
		)
		if err != nil {
			logger.Debugf("error writing service access data to kubernetes (%s)", err)
			http.Error(w, "error writing service access data to kubernetes", http.StatusInternalServerError)
			return
		}
		fullResp := fullBindingResponse{
			ConfigMapAndSecret: configMapAndSecret{
				ConfigMapInfo: configMapQualified,
				SecretsInfo:   secretsQualified,
			},
		}
		// TODO: write creds to a ConfigMap and Secrets as necessary, and return information on where the map and secrets are located
		if err := json.NewEncoder(w).Encode(fullResp); err != nil {
			logger.Debugf("error encoding response to client (%s)", err)
			http.Error(w, "error encoding JSON response", http.StatusInternalServerError)
			return
		}
	})
}
