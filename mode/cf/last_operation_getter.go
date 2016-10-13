package cf

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/deis/steward/mode"
)

// LastOperationGetter fetches the last operation performed after an async provision or deprovision response
type lastOperationGetter struct {
	ctx     context.Context
	cl      *restClient
	timeout time.Duration
}

func (l *lastOperationGetter) GetLastOperation(
	serviceID,
	planID,
	operation,
	instanceID string,
) (*mode.GetLastOperationResponse, error) {
	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, serviceID)
	query.Add(planIDQueryKey, planID)
	query.Add(operationQueryKey, operation)
	req, err := l.cl.Get(query, "v2", "service_instances", instanceID, "last_operation")
	if err != nil {
		return nil, err
	}
	ctx, cancelFn := context.WithTimeout(l.ctx, l.timeout)
	defer cancelFn()
	res, err := l.cl.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// An HTTP response code of 410 (gone) is a distinct state that deprovision may wish to
	// interpret as success.
	if res.StatusCode == http.StatusGone {
		return &mode.GetLastOperationResponse{
			State: mode.LastOperationStateGone.String(),
		}, nil
	}
	resp := new(mode.GetLastOperationResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func newLastOperationGetter(ctx context.Context, cl *restClient, callTimeout time.Duration) mode.LastOperationGetter {
	return &lastOperationGetter{
		ctx:     ctx,
		cl:      cl,
		timeout: callTimeout,
	}
}
