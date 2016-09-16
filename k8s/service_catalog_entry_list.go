package k8s

import (
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/unversioned"
)

// ServiceCatalogEntryList is the third party resource that represents a list of service catalog entries
type ServiceCatalogEntryList struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Items                []*ServiceCatalogEntry `json:"items"`
}
