package k8s

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/juju/loggo"
)

func init() {
	logger.SetLogLevel(loggo.TRACE)
}

func TestResourceAPIVersion(t *testing.T) {
	const (
		ver = "testver"
	)
	assert.Equal(t, resourceAPIVersion(ver), fmt.Sprintf("%s/%s", resourceAPIVersionBase, ver), "api version")
}
