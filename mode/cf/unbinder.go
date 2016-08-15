package cf

import (
	"net/http"
	"net/url"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

const (
	serviceIDQueryKey = "service_id"
	planIDQueryKey    = "plan_id"
)

type unbinder struct {
	cl *RESTClient
}

func (u unbinder) Unbind(serviceID, planID, instanceID, bindingID string) error {
	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, serviceID)
	query.Add(planIDQueryKey, planID)
	req, err := u.cl.Delete(query, "v2", "service_instances", instanceID, "service_bindings", bindingID)
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
func NewUnbinder(cl *RESTClient) mode.Unbinder {
	return unbinder{cl: cl}
}
