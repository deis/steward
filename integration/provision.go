package main

import (
	"github.com/deis/steward/mode"
	"github.com/juju/loggo"
	"github.com/pborman/uuid"
)

func provision(
	logger loggo.Logger,
	provisioner mode.Provisioner,
	svcID,
	planID,
	instanceID string,
) (*mode.ProvisionRequest, *mode.ProvisionResponse, error) {
	req := &mode.ProvisionRequest{
		OrganizationGUID: uuid.New(),
		PlanID:           planID,
		ServiceID:        svcID,
		SpaceGUID:        uuid.New(),
	}
	resp, err := provisioner.Provision(instanceID, req)
	if err != nil {
		return nil, nil, err
	}
	return req, resp, nil
}
