package web

import (
	"net/http"
)

// WithBasicAuth is HTTP middleware that executes next if the incoming request has basic auth equivalent to the auth described in creds
func WithBasicAuth(creds *BasicAuth, next http.Handler) http.Handler {
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
