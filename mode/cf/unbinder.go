package cf

import (
	"net/http"
	"net/url"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/juju/loggo"
)

const (
	serviceIDQueryKey = "service_id"
	planIDQueryKey    = "plan_id"
)

type unbinder struct {
	logger loggo.Logger
	cl     *RESTClient
}

func (u unbinder) Unbind(serviceID, planID, instanceID, bindingID string) error {
	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, serviceID)
	query.Add(planIDQueryKey, planID)
	req, err := u.cl.Delete(u.logger, query, "v2", "service_instances", instanceID, "service_bindings", bindingID)
	if err != nil {
		return err
	}
	resp, err := u.cl.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return web.ErrUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: resp.StatusCode}
	}
	return nil
}

// NewUnbinder returns a CloudFoundry implementation of a mode.Unbinder
func NewUnbinder(logger loggo.Logger, cl *RESTClient) mode.Unbinder {
	return unbinder{logger: logger, cl: cl}
}
