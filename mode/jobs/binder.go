package jobs

import (
	"encoding/json"
	"fmt"

	"github.com/deis/steward/mode"
)

type binder struct {
	pr *podRunner
}

func (b binder) Bind(instanceID, bindingID string, bindRequest *mode.BindRequest) (*mode.BindResponse, error) {
	podName := fmt.Sprintf("bind-%s", bindingID)
	output, err := b.pr.run(podName, brokerBinary, "bind",
		"--instance-id", instanceID,
		"--binding-id", bindingID,
		"--service-id", bindRequest.ServiceID,
		"--plan-id", bindRequest.PlanID,
	)
	if err != nil {
		return nil, err
	}
	resp := &mode.BindResponse{}
	if err := json.Unmarshal([]byte(output), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// newBinder creates a new jobs-broker-backed binder implementation
func newBinder(pr *podRunner) mode.Binder {
	return binder{pr: pr}
}
