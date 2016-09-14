package helm

import (
	"testing"

	"github.com/arschles/assert"
)

func TestProvisionBehaviorFromString(t *testing.T) {
	b1, err := provisionBehaviorFromString("noop")
	assert.NoErr(t, err)
	assert.Equal(t, b1, ProvisionBehaviorNoop, "returned provision behavior")
	b2, err := provisionBehaviorFromString("active")
	assert.NoErr(t, err)
	assert.Equal(t, b2, ProvisionBehaviorActive, "returned provision behavior")
	b3, err := provisionBehaviorFromString("unknown")
	assert.Equal(t, b3, ProvisionBehavior(""), "returned provision behavior")
	if _, ok := err.(errUnknownProvisionBehavior); !ok {
		t.Fatalf("returned error was not an errUnknownProvisionBehavior")
	}
}
