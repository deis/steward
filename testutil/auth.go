package testutil

import (
	"github.com/deis/steward/web"
)

// GetAuth returns a basic authentication struct
func GetAuth() *web.BasicAuth {
	return &web.BasicAuth{Username: "testuser", Password: "testpass"}
}
