package common

import (
	"github.com/deis/steward/mode"
)

// NoopLastOperationGetter is a mode.LastOperationGetter that always returns a successful response
type NoopLastOperationGetter struct{}

// GetLastOperation is the LastOperationGetter interface implementation. It always returns a successful response
func (n NoopLastOperationGetter) GetLastOperation(instanceID string) (*mode.GetLastOperationResponse, error) {
	return &mode.GetLastOperationResponse{State: mode.LastOperationStateSucceeded.String()}, nil
}
