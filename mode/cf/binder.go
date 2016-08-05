package cf

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/juju/loggo"
)

type bindRequest struct {
	ServiceID  string          `json:"service_id"`
	PlanID     string          `json:"plan_id"`
	Parameters mode.JSONObject `json:"parameters"`
}

type bindResponse struct {
	Credentials mode.JSONObject `json:"parameters"`
}

type binder struct {
	logger loggo.Logger
	cl     *RESTClient
}

func (b binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	bodyBytes := new(bytes.Buffer)
	if err := json.NewDecoder(bodyBytes).Decode(bindRequest); err != nil {
		return nil, err
	}

	req, err := b.cl.Put(b.logger, emptyQuery, bodyBytes, "v2", "service_instances", instanceID, "service_bindings", bindingID)
	if err != nil {
		return nil, err
	}

	res, err := b.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusConflict || res.StatusCode == web.StatusUnprocessableEntity {
		return nil, web.ErrUnexpectedResponseCode{URL: req.URL.String(), Expected: http.StatusOK, Actual: res.StatusCode}
	}

	resp := new(bindResponse)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, err
	}
	publicCreds := map[string]string{}
	privateCreds := map[string]string{}
	for credKey, credVal := range resp.Credentials {
		if credKey == "password" {
			privateCreds[credKey] = credVal
		} else {
			publicCreds[credKey] = credVal
		}
	}
	return &mode.BindResponse{
		Status:       res.StatusCode,
		PublicCreds:  publicCreds,
		PrivateCreds: privateCreds,
	}, nil
}

// NewBinder creates a new CloudFoundry-broker-backed binder implementation
func NewBinder(logger loggo.Logger, cl *RESTClient) mode.Binder {
	return binder{logger: logger, cl: cl}
}
