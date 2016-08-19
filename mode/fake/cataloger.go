package fake

import (
	"github.com/deis/steward/mode"
)

// Cataloger is a fake (github.com/deis/steward/mode).Cataloger implementation, suitable for use in unit tests
type Cataloger struct {
	Services []*mode.Service
}

// List is the Cataloger interface implementation. It returns f.Services
func (f Cataloger) List() ([]*mode.Service, error) {
	return f.Services, nil
}
