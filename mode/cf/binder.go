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

type binder struct {
	cl          *RESTClient
	callTimeout time.Duration
	baseCtx     context.Context
}

func (b binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	bodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(bodyBytes).Encode(bindRequest); err != nil {
		return nil, err
	}

	req, err := b.cl.Put(emptyQuery, bodyBytes, "v2", "service_instances", instanceID, "service_bindings", bindingID)
	if err != nil {
		return nil, err
	}

	ctx, cancelFn := context.WithTimeout(b.baseCtx, b.callTimeout)
	defer cancelFn()
	res, err := b.cl.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, web.ErrUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: res.StatusCode}
	}

	resp := new(mode.BindResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	logger.Debugf("got response %+v from backing broker", *resp)
	return resp, nil
}

// NewBinder creates a new CloudFoundry-broker-backed binder implementation
func NewBinder(baseCtx context.Context, cl *RESTClient, callTimeout time.Duration) mode.Binder {
	return binder{cl: cl, callTimeout: callTimeout, baseCtx: baseCtx}
}
