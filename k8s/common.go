package k8s

import (
	"fmt"

	"github.com/juju/loggo"
)

const (
	// TODO: consider making this configurable
	resourceAPIVersionBase = "steward.deis.com"
	apiVersionV1           = "v1"
)

var (
	logger = loggo.GetLogger("k8s")
)

func resourceAPIVersion(v string) string {
	return fmt.Sprintf("%s/%s", resourceAPIVersionBase, v)
}
