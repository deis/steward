package cf

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/juju/loggo"
)

type backendProvisionResp struct {
	Operation string `json:"operation"`
}

type provisioner struct {
	logger loggo.Logger
	cl     *RESTClient
}

func (p provisioner) Provision(instanceID string, pReq *mode.ProvisionRequest) (*mode.ProvisionResponse, error) {
	bodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBytes).Encode(pReq); err != nil {
		return nil, err
	}
	req, err := p.cl.Put(p.logger, emptyQuery, bodyBytes, "v2", "service_instances", instanceID)
	if err != nil {
		return nil, err
	}
	res, err := p.cl.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusConflict && res.StatusCode == web.StatusUnprocessableEntity {
		return nil, web.ErrUnexpectedResponseCode{
			URL:      req.URL.String(),
			Expected: http.StatusOK,
			Actual:   res.StatusCode,
		}
	}
	resp := new(backendProvisionResp)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return &mode.ProvisionResponse{Status: res.StatusCode, Operation: resp.Operation}, nil
}

// NewProvisioner creates a new CloudFoundry-broker-backed provisioner implementation
func NewProvisioner(logger loggo.Logger, cl *RESTClient) mode.Provisioner {
	return provisioner{logger: logger, cl: cl}
}
