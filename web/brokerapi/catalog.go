package brokerapi

import (
	"encoding/json"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/juju/loggo"
)

type catalogResp struct {
	Services []*mode.Service `json:"services"`
}

func catalogHandler(logger loggo.Logger, cataloger mode.Cataloger, auth *web.BasicAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "basic auth must be provided", http.StatusUnauthorized)
			return
		}
		if username != auth.Username || password != auth.Password {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		svcs, err := cataloger.List()
		if err != nil {
			logger.Debugf("error listing services (%s)", err)
			http.Error(w, "error listing services", http.StatusInternalServerError)
			return
		}
		resp := catalogResp{Services: svcs}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Debugf("error encoding response JSON (%s)", err)
			http.Error(w, "error encoding response JSON", http.StatusInternalServerError)
			return
		}
	})
}
