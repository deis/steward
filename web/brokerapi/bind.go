package brokerapi

import (
	"encoding/json"
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/gorilla/mux"
)

type bindResponse struct {
	Credentials mode.JSONObject `json:"credentials"`
}

func bindingHandler(binder mode.Binder, cmCreator k8s.ConfigMapCreator) http.Handler {
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

		bindReq := new(mode.BindRequest)
		if err := json.NewDecoder(r.Body).Decode(bindReq); err != nil {
			logger.Debugf("error decoding bind request (%s)", err)
			http.Error(w, "error decoding bind request", http.StatusBadRequest)
			return
		}

		targetNamespace, err := bindReq.TargetNamespace()
		if err != nil {
			logger.Debugf("target namespace is missing (%s)", err)
			http.Error(w, "target namespace param is missing", http.StatusBadRequest)
			return
		}

		targetName, err := bindReq.TargetName()
		if err != nil {
			logger.Debugf("target name is missing (%s)", err)
			http.Error(w, "target name param is missing", http.StatusBadRequest)
		}

		bindRes, err := binder.Bind(instanceID, bindingID, bindReq)
		if err != nil {
			logger.Debugf("error binding (%s)", err)
			http.Error(w, "error binding", http.StatusInternalServerError)
			return
		}

		logger.Debugf("writing creds %+v to configmap %s/%s", bindRes.Creds, targetNamespace, targetName)

		if err := writeToKubernetes(targetNamespace, targetName, bindRes.Creds, cmCreator); err != nil {
			logger.Debugf("error writing service access data to kubernetes (%s)", err)
			http.Error(w, "error writing service access data to kubernetes", http.StatusInternalServerError)
			return
		}
		fullResp := bindResponse{
			Credentials: mode.JSONObject(map[string]string{
				mode.TargetNameKey:      targetName,
				mode.TargetNamespaceKey: targetNamespaceKey,
			}),
		}
		if err := json.NewEncoder(w).Encode(fullResp); err != nil {
			logger.Debugf("error encoding response to client (%s)", err)
			http.Error(w, "error encoding JSON response", http.StatusInternalServerError)
			return
		}
	})
}
