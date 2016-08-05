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
	auth *web.BasicAuth,
	configMapDeleter k8s.ConfigMapDeleter,
	secretDeleter k8s.SecretDeleter,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "missing basic auth", http.StatusUnauthorized)
			return
		}
		if username != auth.Username || password != auth.Password {
			http.Error(w, "username or password are invalid", http.StatusUnauthorized)
			return
		}
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
			logger.Debugf("error unbinding (%s)", err)
			http.Error(w, "error unbinding", http.StatusInternalServerError)
			return
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
