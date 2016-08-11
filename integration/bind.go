package main

import (
	"bytes"
	"encoding/json"
	"net/url"

	"github.com/deis/steward/mode/cf"
	"github.com/juju/loggo"
)

type bindReq struct {
	ServiceID  string            `json:"service_id"`
	PlanID     string            `json:"plan_id"`
	Parameters map[string]string `json:"parameters"`
}

type qualifiedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type bindResp struct {
	ConfigMapInfo qualifiedName   `json:"config_map_info"`
	SecretsInfo   []qualifiedName `json:"secrets_info"`
}

func bind(
	logger loggo.Logger,
	cl *cf.RESTClient,
	svcID,
	planID,
	instanceID,
	bindID,
	targetNS,
	targetName string,
) (*bindResp, error) {
	bindReq := &bindReq{
		ServiceID: svcID,
		PlanID:    planID,
		Parameters: map[string]string{
			"target_namespace": targetNS,
			"target_name":      targetName,
		},
	}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(bindReq); err != nil {
		return nil, err
	}

	req, err := cl.Put(
		logger,
		url.Values(map[string][]string{}),
		reqBody,
		"v2",
		"service_instances",
		instanceID,
		"service_bindings",
		bindID,
	)
	if err != nil {
		return nil, err
	}

	rawResp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawResp.Body.Close()

	resp := new(bindResp)
	if err := json.NewDecoder(rawResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
