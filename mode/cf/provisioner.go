package cf

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

type provisioner struct {
	cl          *restClient
	baseCtx     context.Context
	callTimeout time.Duration
}

func (p provisioner) Provision(instanceID string, pReq *mode.ProvisionRequest) (*mode.ProvisionResponse, error) {
	query := url.Values(map[string][]string{})
	query.Add(asyncQueryKey, "true")
	bodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBytes).Encode(pReq); err != nil {
		return nil, err
	}
	req, err := p.cl.Put(query, bodyBytes, "v2", "service_instances", instanceID)
	if err != nil {
		return nil, err
	}
	ctx, cancelFn := context.WithTimeout(p.baseCtx, p.callTimeout)
	defer cancelFn()
	res, err := p.cl.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := new(mode.ProvisionResponse)
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		resp.IsAsync = false
	case http.StatusAccepted:
		resp.IsAsync = true
	default:
		return nil, web.ErrUnexpectedResponseCode{
			URL:      req.URL.String(),
			Expected: http.StatusOK,
			Actual:   res.StatusCode,
		}
	}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// newProvisioner creates a new CloudFoundry-broker-backed provisioner implementation
func newProvisioner(baseCtx context.Context, cl *restClient, callTimeout time.Duration) mode.Provisioner {
	return provisioner{cl: cl, baseCtx: baseCtx, callTimeout: callTimeout}
}
