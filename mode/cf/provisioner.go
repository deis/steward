package cf

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

type provisioner struct {
	cl          *RESTClient
	baseCtx     context.Context
	callTimeout time.Duration
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
	ctx, cancelFn := context.WithTimeout(p.baseCtx, p.callTimeout)
	defer cancelFn()
	res, err := p.cl.Do(ctx, req)
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
	resp := new(mode.ProvisionResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// NewProvisioner creates a new CloudFoundry-broker-backed provisioner implementation
func NewProvisioner(baseCtx context.Context, cl *RESTClient, callTimeout time.Duration) mode.Provisioner {
	return provisioner{cl: cl, baseCtx: baseCtx, callTimeout: callTimeout}
}
