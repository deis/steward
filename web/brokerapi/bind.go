package brokerapi

import (
	"encoding/json"
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

type fullBindingResponse struct {
	ConfigMapAndSecret configMapAndSecret `json:"credentials"`
}

// a pair of qualified names describing the config map and associated secrets
type configMapAndSecret struct {
	ConfigMapInfo *qualifiedName   `json:"config_map_info"`
	SecretsInfo   []*qualifiedName `json:"secrets_info"`
}

// a (name, namespace) pair to identify exactly where a resource is
type qualifiedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func bindingHandler(
	logger loggo.Logger,
	binder mode.Binder,
	auth *web.BasicAuth,
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
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authorization missing", http.StatusBadRequest)
			return
		}
		if username != auth.Username || password != auth.Password {
			http.Error(w, "wrong login credentials", http.StatusUnauthorized)
			return
		}

		bindReq := new(mode.BindRequest)
		if err := json.NewDecoder(r.Body).Decode(bindReq); err != nil {
			logger.Debugf("error decoding bind request (%s)", err)
			http.Error(w, "error decoding bind request", http.StatusBadRequest)
			return
		}

		targetNamespace, err := bindReq.Parameters.String(targetNamespaceKey)
		if err != nil {
			logger.Debugf("parameter %s is missing or malformed (%s)", targetNamespaceKey, err)
			http.Error(w, "target namespace param is missing or malformed", http.StatusBadRequest)
			return
		}

		bindRes, err := binder.Bind(instanceID, bindingID, bindReq)
		if err != nil {
			logger.Debugf("error binding (%s)", err)
			http.Error(w, "error binding", http.StatusInternalServerError)
			return
		}

		configMapQualified, secretsQualified, err := writeToKubernetes(
			bindReq.ServiceID,
			bindReq.PlanID,
			targetNamespace,
			bindRes.PublicCreds,
			bindRes.PrivateCreds,
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
		if err := json.NewEncoder(w).Encode(fullResp); err != nil {
			logger.Debugf("error encoding response to client (%s)", err)
			http.Error(w, "error encoding JSON response", http.StatusInternalServerError)
			return
		}
	})
}
