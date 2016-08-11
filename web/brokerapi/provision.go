package brokerapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func provisioningHandler(logger loggo.Logger, provisioner mode.Provisioner) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceID, ok := vars[instanceIDPathKey]
		if !ok {
			http.Error(w, "missing instance ID", http.StatusBadRequest)
			return
		}
		provisionReq := new(mode.ProvisionRequest)
		if err := json.NewDecoder(r.Body).Decode(provisionReq); err != nil {
			logger.Debugf("error decoding provision request (%s)", err)
			http.Error(w, "error decoding provision request", http.StatusBadRequest)
			return
		}
		resp, err := provisioner.Provision(instanceID, provisionReq)
		if err != nil {
			switch t := err.(type) {
			case web.ErrUnexpectedResponseCode:
				logger.Debugf("expected response code %d, got %d for provisioning (%s)", t.Expected, t.Actual, t.Expected)
				http.Error(w, "error provisioning. backend returned failure response", t.Actual)
				return
			default:
				logger.Debugf("error provisioning (%s)", err)
				http.Error(w, "error provisioning", http.StatusInternalServerError)
				return
			}
		}
		respStr := fmt.Sprintf(`{"operation":"%s"}`, resp.Operation)
		w.WriteHeader(resp.Status)
		w.Write([]byte(respStr))
	})
}
