package cf

import (
	"fmt"
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
	Scheme   string
	Host     string
	Username string
	Password string
}

// NewClient creates a new CloudFoundry client with the given HTTP client, host, username and password
func NewClient(cl *http.Client, scheme, host, user, pass string) *Client {
	return &Client{Client: cl, Scheme: scheme, Host: host, Username: user, Password: pass}
}

// Get creates a GET request with the given path
func (c *Client) Get(logger loggo.Logger, pathElts ...string) (*http.Request, error) {
	logger, _ = loggo.NewLogger("cf", logger)
	pathStr := strings.Join(pathElts, "/")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s:%s@%s/%s", c.Scheme, c.Username, c.Password, c.Host, pathStr), nil)
	if err != nil {
		logger.Debugf("CF Client Get error (%s)", err)
		return nil, err
	}
	logger.Debugf("CF client making request to %s", req.URL.String())
	req.Header.Set(versionHeader, apiVersion)
	return req, nil
}
