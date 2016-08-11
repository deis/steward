package brokerapi

import (
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func deprovisioningHandler(logger loggo.Logger, deprovisioner mode.Deprovisioner) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceID, ok := vars[instanceIDPathKey]
		if !ok {
			http.Error(w, "missing instance ID in path", http.StatusBadRequest)
			return
		}
		serviceID := r.URL.Query().Get(serviceIDQueryKey)
		planID := r.URL.Query().Get(planIDQueryKey)
		if serviceID == "" {
			http.Error(w, "missing service ID in query string", http.StatusBadRequest)
			return
		}
		if planID == "" {
			http.Error(w, "missing plan ID in query string", http.StatusBadRequest)
			return
		}
		resp, err := deprovisioner.Deprovision(instanceID, serviceID, planID)
		if err != nil {
			switch t := err.(type) {
			case web.ErrUnexpectedResponseCode:
				logger.Debugf("expected response code %d, got %d for deprovision (%s)", t.Expected, t.Actual, t.Expected)
				http.Error(w, "error deprovisioning. backend returned failure response", t.Actual)
				return
			default:
				logger.Debugf("error deprovisioning (%s)", err)
				http.Error(w, "error deprovisioning", http.StatusInternalServerError)
				return
			}
		}

		respStr := fmt.Sprintf(`{"operation":"%s"}`, resp.Operation)
		w.WriteHeader(resp.Status)
		w.Write([]byte(respStr))
	})
}
