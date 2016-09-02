package helm

import (
	"fmt"
)

type errUnknownProvisionBehavior struct {
	behavior string
}

func (e errUnknownProvisionBehavior) Error() string {
	return fmt.Sprintf("unknown provision behavior %s", e.behavior)
}

const (
	// ProvisionBehaviorNoop indicates that steward should 'helm install' a new instance of the chart on startup, and that provision and deprovisoin should do nothing
	ProvisionBehaviorNoop ProvisionBehavior = "noop"
	// ProvisionBehaviorActive indicates that steward should 'helm install' a new instance of the chart on every provision operation (and helm uninstall on each deprovision operation)
	ProvisionBehaviorActive ProvisionBehavior = "active"
)

// ProvisionBehavior is the indication for what steward should do in helm mode when a provision comes in. It implements fmt.Stringer
type ProvisionBehavior string

// ProvisionBehaviorFromString returns the ProvisionBehavior that corresponds to s. If s is an invalid provision behavior, returns an empty string and a non-nil error
func ProvisionBehaviorFromString(s string) (ProvisionBehavior, error) {
	switch s {
	case ProvisionBehaviorActive.String():
		return ProvisionBehaviorActive, nil
	case ProvisionBehaviorNoop.String():
		return ProvisionBehaviorNoop, nil
	default:
		return "", errUnknownProvisionBehavior{behavior: s}
	}
}

func (p ProvisionBehavior) String() string {
	return string(p)
}
