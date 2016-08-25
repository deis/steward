package jobs

import (
	"fmt"

	"github.com/deis/steward/mode"
)

type unbinder struct {
	pr *podRunner
}

func (u unbinder) Unbind(instanceID, bindingID string, uReq *mode.UnbindRequest) error {
	podName := fmt.Sprintf("unbind-%s", bindingID)
	_, err := u.pr.run(podName, brokerBinary, "unbind",
		"--instance-id", instanceID,
		"--binding-id", bindingID,
		"--service-id", uReq.ServiceID,
		"--plan-id", uReq.PlanID,
	)
	if err != nil {
		return err
	}
	return nil
}

// newUnbinder returns a jobs implementation of a mode.Unbinder
func newUnbinder(pr *podRunner) mode.Unbinder {
	return unbinder{pr: pr}
}
