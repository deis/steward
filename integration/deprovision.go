package main

import (
	"github.com/deis/steward/mode"
	"github.com/juju/loggo"
)

func deprovision(
	logger loggo.Logger,
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
