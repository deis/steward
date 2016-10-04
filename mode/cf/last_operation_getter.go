package cf

import (
	"github.com/deis/steward/mode"
)

// LastOperationGetter fetches the last operation performed after an async provision or deprovision response
type lastOperationGetter struct{}

func (l *lastOperationGetter) GetLastOperation(instanceID string) (*mode.GetLastOperationResponse, error) {
	return nil, nil
}
