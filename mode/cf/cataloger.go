package cf

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

type cataloger struct {
	cl          *restClient
	baseCtx     context.Context
	callTimeout time.Duration
}

func (c cataloger) List() ([]*mode.Service, error) {
	req, err := c.cl.Get(emptyQuery, "v2", "catalog")
	if err != nil {
		return nil, err
	}
	ctx, cancelFn := context.WithTimeout(c.baseCtx, c.callTimeout)
	defer cancelFn()
	res, err := c.cl.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	serviceList := new(mode.ServiceList)
	// TODO: drain the response body to avoid a connection leak
	if err := json.NewDecoder(res.Body).Decode(serviceList); err != nil {
		logger.Debugf("error decoding JSON response body from backend CF broker (%s)", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, web.ErrUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: res.StatusCode}
	}
	return serviceList.Services, nil
}

// newCataloger returns a new Cataloger implementation, backed by a CF service broker
func newCataloger(baseCtx context.Context, cl *restClient, callTimeout time.Duration) mode.Cataloger {
	return cataloger{cl: cl, baseCtx: baseCtx, callTimeout: callTimeout}
}
