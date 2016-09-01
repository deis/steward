package cf

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
)

type binder struct {
	cl *RESTClient
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

	res, err := b.cl.Do(req)
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
func NewBinder(cl *RESTClient) mode.Binder {
	return binder{cl: cl}
}
