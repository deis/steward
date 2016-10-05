package cf

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

const (
	serviceIDQueryKey = "service_id"
	planIDQueryKey    = "plan_id"
	operationQueryKey = "operation"
)

type unbinder struct {
	cl          *restClient
	baseCtx     context.Context
	callTimeout time.Duration
}

func (u unbinder) Unbind(instanceID, bindingID string, uReq *mode.UnbindRequest) error {
	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, uReq.ServiceID)
	query.Add(planIDQueryKey, uReq.PlanID)
	req, err := u.cl.Delete(query, "v2", "service_instances", instanceID, "service_bindings", bindingID)
	if err != nil {
		return err
	}
	ctx, cancelFn := context.WithTimeout(u.baseCtx, u.callTimeout)
	defer cancelFn()
	resp, err := u.cl.Do(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return web.ErrUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: resp.StatusCode}
	}
	return nil
}

// newUnbinder returns a CloudFoundry implementation of a mode.Unbinder
func newUnbinder(baseCtx context.Context, cl *restClient, callTimeout time.Duration) mode.Unbinder {
	return unbinder{cl: cl, baseCtx: baseCtx, callTimeout: callTimeout}
}
