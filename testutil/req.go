package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/arschles/testsrv"
	"github.com/deis/steward/web"
)

// NewReq creates a new http.Request to be executed against srv, with the given auth, HTTP method, query string, body, and path
func NewReq(
	srv *testsrv.Server,
	auth *web.BasicAuth,
	method string,
	query url.Values,
	body interface{},
	path ...string,
) (*http.Request, error) {

	urlStr := fmt.Sprintf("%s/%s", srv.URLStr(), strings.Join(path, "/"))

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	req.SetBasicAuth(auth.Username, auth.Password)
	return req, nil
}
