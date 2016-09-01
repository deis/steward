package cf

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

type deprovisioner struct {
	cl *RESTClient
}

func (d deprovisioner) Deprovision(instanceID string, dReq *mode.DeprovisionRequest) (*mode.DeprovisionResponse, error) {
	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, dReq.ServiceID)
	query.Add(planIDQueryKey, dReq.PlanID)
	req, err := d.cl.Delete(query, "v2", "service_instances", instanceID)
	if err != nil {
		return nil, err
	}
	res, err := d.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return nil, web.ErrUnexpectedResponseCode{
			URL:      req.URL.String(),
			Expected: http.StatusOK,
			Actual:   res.StatusCode,
		}
	}
	resp := new(mode.DeprovisionResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// NewDeprovisioner creates a new CloudFoundry-broker-backed deprovisioner implementation
func NewDeprovisioner(cl *RESTClient) mode.Deprovisioner {
	return deprovisioner{cl: cl}
}
