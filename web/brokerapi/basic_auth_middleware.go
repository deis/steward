package brokerapi

import (
	"net/http"

	"github.com/deis/steward/web"
)

func withBasicAuth(creds *web.BasicAuth, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authorization missing", http.StatusBadRequest)
			return
		}
		if username != creds.Username {
			http.Error(w, "incorrect username", http.StatusUnauthorized)
			return
		}
		if password != creds.Password {
			http.Error(w, "incorrect password", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
