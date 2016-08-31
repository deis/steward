package cf

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

type backendProvisionResp struct {
	Operation string `json:"operation"`
}

type provisioner struct {
	cl *RESTClient
}

func (p provisioner) Provision(instanceID string, pReq *mode.ProvisionRequest) (*mode.ProvisionResponse, error) {
	bodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBytes).Encode(pReq); err != nil {
		return nil, err
	}
	req, err := p.cl.Put(emptyQuery, bodyBytes, "v2", "service_instances", instanceID)
	if err != nil {
		return nil, err
	}
	res, err := p.cl.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
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
	return &mode.ProvisionResponse{Operation: resp.Operation}, nil
}

// NewProvisioner creates a new CloudFoundry-broker-backed provisioner implementation
func NewProvisioner(cl *RESTClient) mode.Provisioner {
	return provisioner{cl: cl}
}
