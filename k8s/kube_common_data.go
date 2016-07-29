package k8s

import (
	"fmt"
)

const (
	// TODO: consider making this configurable
	resourceAPIVersionBase = "steward.deis.com"
)

// KubeCommonData is a json-encodable structure that represents data common to almost all kubernetes resources
type KubeCommonData struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

// KubeCommonMetadata is a json-encodable structure that represents kubernetes common metadata
type KubeCommonMetadata struct {
	Name string `json:"name"`
}

func resourceAPIVersion(v string) string {
	return fmt.Sprintf("%s/%s", resourceAPIVersionBase, v)
}
