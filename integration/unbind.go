package main

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/deis/steward/mode/cf"
	"github.com/juju/loggo"
)

var (
	errUnbindNonEmptyResp = errors.New("non-empty response body returned from unbind")
)

func unbind(
	logger loggo.Logger,
	cl *cf.RESTClient,
	svcID,
	planID,
	instanceID,
	bindID,
	targetNS,
	targetName string,
) error {
	query := url.Values(map[string][]string{
		"target-namespace": []string{targetNS},
		"target-name":      []string{targetName},
		"service_id":       []string{svcID},
		"plan_id":          []string{planID},
	})
	req, err := cl.Delete(logger, query, "v2", "service_instances", svcID, "service_bindings", bindID)
	if err != nil {
		return err
	}
	rawResp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer rawResp.Body.Close()
	resp := make(map[string]string)
	if err := json.NewDecoder(rawResp.Body).Decode(resp); err != nil {
		return nil
	}
	if len(resp) > 0 {
		return errUnbindNonEmptyResp
	}
	return nil
}
