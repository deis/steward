package cf

import (
	"encoding/json"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/juju/loggo"
)

// wrapper for a list of services, returned by the backend broker
type serviceList struct {
	Services []*mode.Service `json:"services"`
}

type cataloger struct {
	logger loggo.Logger
	cl     *RESTClient
}

func (c cataloger) List() ([]*mode.Service, error) {
	req, err := c.cl.Get(c.logger, emptyQuery, "v2", "catalog")
	if err != nil {
		return nil, err
	}
	res, err := c.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	serviceList := new(serviceList)
	// TODO: drain the response body to avoid a connection leak
	if err := json.NewDecoder(res.Body).Decode(serviceList); err != nil {
		c.logger.Debugf("error decoding JSON response body from backend CF broker (%s)", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, web.ErrUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: res.StatusCode}
	}
	return serviceList.Services, nil
}

// NewCataloger returns a new Cataloger implementation, backed by a CF service broker
func NewCataloger(logger loggo.Logger, cl *RESTClient) mode.Cataloger {
	return cataloger{logger: logger, cl: cl}
}
