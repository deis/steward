package fake

import (
	"github.com/deis/steward/mode"
)

// GetLastOperationCall represents the parameters of a call to GetLastOperation
type GetLastOperationCall struct {
	ServiceID  string
	PlanID     string
	Operation  string
	InstanceID string
}

// LastOperationGetter is a fake implementation of a mode.LastOperationGetter. It's useful for use in unit tests
type LastOperationGetter struct {
	Calls []*GetLastOperationCall
	Ret   func() *mode.GetLastOperationResponse
}

// GetLastOperation appends to l.Calls and returns l.Ret(), nil. Not concurrency safe
func (l *LastOperationGetter) GetLastOperation(
	serviceID,
	planID,
	operation,
	instanceID string,
) (*mode.GetLastOperationResponse, error) {

	newCall := &GetLastOperationCall{
		ServiceID:  serviceID,
		PlanID:     planID,
		Operation:  operation,
		InstanceID: instanceID,
	}
	l.Calls = append(l.Calls, newCall)
	return l.Ret(), nil
}
