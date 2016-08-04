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

// Client represents a client to talk to the CloudFoundry API at Host
type Client struct {
	Client   *http.Client
	scheme   string
	host     string
	port     int
	username string
	password string
}

// NewClient creates a new CloudFoundry client with the given HTTP client, host, username and password
func NewClient(cl *http.Client, scheme, host string, port int, user, pass string) *Client {
	return &Client{Client: cl, scheme: scheme, host: host, port: port, username: user, password: pass}
}

// returns the full URL to the broker, including basic auth
func (c Client) fullBaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d", c.scheme, c.username, c.password, c.host, c.port)
}

// returns a fully formed URL string including a path comprised of pathElts
func (c Client) urlStr(pathElts ...string) string {
	pathStr := strings.Join(pathElts, "/")
	return fmt.Sprintf("%s/%s", c.fullBaseURL(), pathStr)
}

// Get creates a GET request with the given path
func (c *Client) Get(logger loggo.Logger, pathElts ...string) (*http.Request, error) {
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
func (c *Client) Put(logger loggo.Logger, body io.Reader, pathElts ...string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.urlStr(pathElts...), body)
	if err != nil {
		logger.Debugf("CF Client PUT error (%s)", err)
		return nil, err
	}
	logger.Debugf("CF client making request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}

// DoPut creates a PUT request with the given path and body, then executes the request using c.Client
func (c *Client) DoPut(logger loggo.Logger, body io.Reader, pathElts ...string) (*http.Response, error) {
	req, err := c.Put(logger, body, pathElts...)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}
