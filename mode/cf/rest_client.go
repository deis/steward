package cf

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/juju/loggo"
)

const (
	apiVersion    = "2.9"
	versionHeader = "X-Broker-Api-Version"
)

var (
	emptyQuery = url.Values(map[string][]string{})
)

// RESTClient represents a client to talk to a CloudFoundry broker API at a given location
type RESTClient struct {
	Client   *http.Client
	scheme   string
	host     string
	port     int
	username string
	password string
}

// NewRESTClient creates a new CloudFoundry client with the given HTTP client, host, username and password
func NewRESTClient(cl *http.Client, scheme, host string, port int, user, pass string) *RESTClient {
	return &RESTClient{Client: cl, scheme: scheme, host: host, port: port, username: user, password: pass}
}

// returns the full URL to the broker, including basic auth
func (c RESTClient) fullBaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d", c.scheme, c.username, c.password, c.host, c.port)
}

// returns a fully formed URL string including a path comprised of pathElts
func (c RESTClient) urlStr(pathElts ...string) string {
	pathStr := strings.Join(pathElts, "/")
	return fmt.Sprintf("%s/%s", c.fullBaseURL(), pathStr)
}

// Get creates a GET request with the given query string values and path, or a non-nil error if request creation failed
func (c *RESTClient) Get(logger loggo.Logger, query url.Values, pathElts ...string) (*http.Request, error) {
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
func (c *RESTClient) Put(logger loggo.Logger, query url.Values, body io.Reader, pathElts ...string) (*http.Request, error) {
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
func (c *RESTClient) Delete(logger loggo.Logger, query url.Values, pathElts ...string) (*http.Request, error) {
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
func (c *RESTClient) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}
