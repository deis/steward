package main

import (
	"github.com/deis/steward/mode"
)

func deprovision(
	deprovisioner mode.Deprovisioner,
	svcID,
	planID,
	instanceID string,
) (*mode.DeprovisionResponse, error) {
	resp, err := deprovisioner.Deprovision(instanceID, svcID, planID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
