package brokerapi

import (
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func unbindHandler(
	logger loggo.Logger,
	unbinder mode.Unbinder,
	configMapDeleter k8s.ConfigMapDeleter,
	secretDeleter k8s.SecretDeleter,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceID, ok := vars[instanceIDPathKey]
		if !ok {
			http.Error(w, "missing instance ID in path", http.StatusBadRequest)
			return
		}
		bindingID, ok := vars[bindingIDPathKey]
		if !ok {
			http.Error(w, "missing binding ID in path", http.StatusBadRequest)
			return
		}
		serviceID := r.URL.Query().Get(serviceIDQueryKey)
		if serviceID == "" {
			http.Error(w, "missing service ID in query string", http.StatusBadRequest)
			return
		}
		planID := r.URL.Query().Get(planIDQueryKey)
		if planID == "" {
			http.Error(w, "missing plan ID in query string", http.StatusBadRequest)
			return
		}
		namespace := r.URL.Query().Get(targetNamespaceKey)
		if namespace == "" {
			http.Error(w, "missing namespace in query string", http.StatusBadRequest)
			return
		}
		if err := unbinder.Unbind(serviceID, planID, instanceID, bindingID); err != nil {
			switch t := err.(type) {
			case web.ErrUnexpectedResponseCode:
				logger.Debugf("expected response code %d, got %d for unbinding (%s)", t.Expected, t.Actual, t.Expected)
				http.Error(w, "error unbinding. backend returned failure response", t.Actual)
				return
			default:
				logger.Debugf("error unbinding (%s)", err)
				http.Error(w, "error unbinding", http.StatusInternalServerError)
				return
			}
		}
		if err := deleteFromKubernetes(
			serviceID,
			planID,
			bindingID,
			instanceID,
			namespace,
			configMapDeleter,
			secretDeleter,
		); err != nil {
			logger.Debugf("error deleting bind resources from kubernetes (%s)", err)
			http.Error(w, "error deleting bind resources from kubernetes", http.StatusInternalServerError)
			return
		}
	})
}
