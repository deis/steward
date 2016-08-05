package web

import (
	"fmt"
)

// ErrUnexpectedResponseCode is an error implementation, intended to be returned from an HTTP client that receives an HTTP response code it didn't expect
type ErrUnexpectedResponseCode struct {
	URL      string
	Actual   int
	Expected int
}

// Error is the error interface implementation
func (e ErrUnexpectedResponseCode) Error() string {
	return fmt.Sprintf("%s - expected response code %d, actual %d", e.URL, e.Expected, e.Actual)
}
