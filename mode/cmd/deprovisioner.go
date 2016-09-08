package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/deis/steward/mode"
)

type deprovisioner struct {
	pr *podRunner
}

func (d deprovisioner) Deprovision(instanceID string, dReq *mode.DeprovisionRequest) (*mode.DeprovisionResponse, error) {
	podName := fmt.Sprintf("deprovision-%s", instanceID)
	output, err := d.pr.run(podName, brokerBinary, "deprovision",
		"--instance-id", instanceID,
		"--service-id", dReq.ServiceID,
		"--plan-id", dReq.PlanID,
	)
	if err != nil {
		return nil, err
	}
	resp := &mode.DeprovisionResponse{}
	if err := json.Unmarshal([]byte(output), resp); err != nil {
		return nil, err
	}
	logger.Infof("Operation: %s", resp.Operation)
	return resp, nil
}

// newDeprovisioner creates a new cmd-broker-backed deprovisioner implementation
func newDeprovisioner(pr *podRunner) mode.Deprovisioner {
	return deprovisioner{pr: pr}
}
