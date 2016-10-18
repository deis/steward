package mode

import (
	"fmt"
)

// ErrUnknownLastOperation is an error indicating that a string is not a valid known last operation
type ErrUnknownLastOperation string

// Error is the error interface implementation
func (e ErrUnknownLastOperation) Error() string {
	return fmt.Sprintf("unknown last operation '%s'", string(e))
}

// LastOperationState represents the state returned in a "get last operation" call. This type implements fmt.Stringer
type LastOperationState string

// String is the fmt.Stringer interface implementation
func (l LastOperationState) String() string {
	return string(l)
}

const (
	// LastOperationStateSucceeded is the LastOperationState indicating that the operation has succeeded
	LastOperationStateSucceeded LastOperationState = "succeeded"
	// LastOperationStateFailed is the LastOperationState indicating that the operation has failed
	LastOperationStateFailed LastOperationState = "failed"
	// LastOperationStateInProgress is the LastOperationState indicating that the operation is still in progress
	LastOperationStateInProgress LastOperationState = "in progress"
	// LastOperationStateGone is the LastOperationState indicating that the broker has deleted the
	// instance in question. In the case of async deprovisioning, this is an indicator of success.
	LastOperationStateGone LastOperationState = "gone"
)

// GetLastOperationResponse is the response body from a get last operation call
type GetLastOperationResponse struct {
	State string `json:"state"`
}

// GetState returns the LastOperationState defined in g.State, or an error if g.State is not a valid LastOperationState
func (g *GetLastOperationResponse) GetState() (LastOperationState, error) {
	switch g.State {
	case LastOperationStateSucceeded.String():
		return LastOperationStateSucceeded, nil
	case LastOperationStateFailed.String():
		return LastOperationStateFailed, nil
	case LastOperationStateInProgress.String():
		return LastOperationStateInProgress, nil
	case LastOperationStateGone.String():
		return LastOperationStateGone, nil
	default:
		return LastOperationState(""), ErrUnknownLastOperation(g.State)
	}
}
