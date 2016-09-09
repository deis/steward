package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
)

type cataloger struct {
	pr *podRunner
}

func (c cataloger) List() ([]*mode.Service, error) {
	podName := fmt.Sprintf("catalog-%s", uuid.New())
	output, err := c.pr.run(podName, brokerBinary, "catalog")
	if err != nil {
		return nil, err
	}
	resp := &mode.ServiceList{}
	if err := json.Unmarshal([]byte(output), resp); err != nil {
		return nil, err
	}
	return resp.Services, nil
}

// newCataloger returns a new Cataloger implementation, backed by a cmd service broker
func newCataloger(pr *podRunner) mode.Cataloger {
	return cataloger{pr: pr}
}
