package brokerapi

import (
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/arschles/testsrv"
	"github.com/deis/steward/web"
)

type okHdl struct{}

func (o okHdl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var (
	creds = web.BasicAuth{Username: "testuser", Password: "testpass"}
)

func TestBasicAuth(t *testing.T) {
	srv := testsrv.StartServer(withBasicAuth(&creds, okHdl{}))
	defer srv.Close()

	resp, err := http.Get(srv.URLStr())
	assert.NoErr(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "response code")

	req, err := http.NewRequest("GET", srv.URLStr(), nil)
	assert.NoErr(t, err)
	req.SetBasicAuth("notausername", creds.Password)
	resp, err = http.DefaultClient.Do(req)
	assert.NoErr(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized, "response code")

	req, err = http.NewRequest("GET", srv.URLStr(), nil)
	assert.NoErr(t, err)
	req.SetBasicAuth(creds.Username, "notapassword")
	resp, err = http.DefaultClient.Do(req)
	assert.NoErr(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized, "response code")

	req, err = http.NewRequest("GET", srv.URLStr(), nil)
	assert.NoErr(t, err)
	req.SetBasicAuth(creds.Username, creds.Password)
	resp, err = http.DefaultClient.Do(req)
	assert.NoErr(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
}
