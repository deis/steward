package cf

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/juju/loggo"
)

const (
	apiVersion    = "2.9"
	versionHeader = "X-Broker-Api-Version"
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

// Get creates a GET request with the given path
func (c *RESTClient) Get(logger loggo.Logger, pathElts ...string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.urlStr(pathElts...), nil)
	if err != nil {
		logger.Debugf("CF Client GET error (%s)", err)
		return nil, err
	}
	logger.Debugf("CF client making request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}

// Put creates a PUT request with the given path and body
func (c *RESTClient) Put(logger loggo.Logger, body io.Reader, pathElts ...string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.urlStr(pathElts...), body)
	if err != nil {
		logger.Debugf("CF Client PUT error (%s)", err)
		return nil, err
	}
	logger.Debugf("CF client making request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}

// Do is a convenience function for c.Client.Do(req)
func (c *RESTClient) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

// DoPut creates a PUT request with the given path and body, then executes the request using c.Client
func (c *RESTClient) DoPut(logger loggo.Logger, body io.Reader, pathElts ...string) (*http.Response, error) {
	req, err := c.Put(logger, body, pathElts...)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}
