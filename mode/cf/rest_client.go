package cf

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/deis/steward/web/ctxhttp"
)

const (
	apiVersion        = "2.9"
	versionHeader     = "X-Broker-Api-Version"
	serviceIDQueryKey = "service_id"
	planIDQueryKey    = "plan_id"
	operationQueryKey = "operation"
	asyncQueryKey     = "accepts_incomplete"
)

var (
	emptyQuery = url.Values(map[string][]string{})
)

// restClient represents a client to talk to a CloudFoundry broker API at a given location
type restClient struct {
	Client   *http.Client
	scheme   string
	host     string
	port     int
	username string
	password string
}

// newRESTClient creates a new CloudFoundry client with the given HTTP client, host, username and password
func newRESTClient(cl *http.Client, scheme, host string, port int, user, pass string) *restClient {
	return &restClient{Client: cl, scheme: scheme, host: host, port: port, username: user, password: pass}
}

// returns the full URL to the broker, including basic auth
func (c restClient) fullBaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d", c.scheme, c.username, c.password, c.host, c.port)
}

// returns a fully formed URL string including a path comprised of pathElts
func (c restClient) urlStr(pathElts ...string) string {
	pathStr := strings.Join(pathElts, "/")
	return fmt.Sprintf("%s/%s", c.fullBaseURL(), pathStr)
}

// Get creates a GET request with the given query string values and path, or a non-nil error if request creation failed
func (c *restClient) Get(query url.Values, pathElts ...string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.urlStr(pathElts...), nil)
	if err != nil {
		logger.Debugf("CF Client GET error (%s)", err)
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	logger.Debugf("CF client making GET request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}

// Put creates a PUT request with the given query string values, request body and path, or a non-nil error if request creation failed
func (c *restClient) Put(query url.Values, body io.Reader, pathElts ...string) (*http.Request, error) {
	req, err := http.NewRequest("PUT", c.urlStr(pathElts...), body)
	if err != nil {
		logger.Debugf("CF Client PUT error (%s)", err)
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	logger.Debugf("CF client making PUT request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}

// Delete creates a DELETE request with the given query string and path, or a non-nil error if request creation failed
func (c *restClient) Delete(query url.Values, pathElts ...string) (*http.Request, error) {
	req, err := http.NewRequest("DELETE", c.urlStr(pathElts...), nil)
	if err != nil {
		logger.Debugf("CF Client DELETE error (%s)", err)
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	logger.Debugf("CF client making DELETE request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}

// Do is a convenience function for c.Client.Do(req)
func (c *restClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return ctxhttp.Do(ctx, c.Client, req)
}
