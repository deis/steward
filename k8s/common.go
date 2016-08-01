package k8s

import (
	"fmt"
)

const (
	// TODO: consider making this configurable
	resourceAPIVersionBase = "steward.deis.com"
	apiVersionV1           = "v1"
)

func resourceAPIVersion(v string) string {
	return fmt.Sprintf("%s/%s", resourceAPIVersionBase, v)
}
