package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/deis/steward/mode"
)

type provisioner struct {
	pr *podRunner
}

func (p provisioner) Provision(instanceID string, provisionRequest *mode.ProvisionRequest) (*mode.ProvisionResponse, error) {
	podName := fmt.Sprintf("provision-%s", instanceID)
	output, err := p.pr.run(podName, brokerBinary, "provision",
		"--instance-id", instanceID,
		"--service-id", provisionRequest.ServiceID,
		"--plan-id", provisionRequest.PlanID,
	)
	if err != nil {
		return nil, err
	}
	resp := &mode.ProvisionResponse{}
	if err := json.Unmarshal([]byte(output), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// newProvisioner creates a new cmd-broker-backed provisioner implementation
func newProvisioner(pr *podRunner) mode.Provisioner {
	return provisioner{pr: pr}
}
