package k8s

import (
	"fmt"

	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	// ServiceCatalogEntryKind is the Kubernetes kind that should be used when creating/updating service catalog entries
	ServiceCatalogEntryKind = "ServiceCatalogEntry"
)

// ServiceCatalogEntry is the third party resource that represents a single service provider + plan. A new ServiceCatalogEntry should be created with NewServiceCatalogEntry
type ServiceCatalogEntry struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Info                 mode.ServiceInfo `json:"service_info"`
	Plan                 mode.ServicePlan `json:"service_plan"`
}

// NewServiceCatalogEntry creates a new ServiceCatalogEntry suitable for writing to the Kubernetes API
func NewServiceCatalogEntry(
	objectMeta api.ObjectMeta,
	info mode.ServiceInfo,
	plan mode.ServicePlan) *ServiceCatalogEntry {
	typeMeta := unversioned.TypeMeta{
		Kind:       ServiceCatalogEntryKind,
		APIVersion: resourceAPIVersion(apiVersionV1),
	}
	objectMeta.Name = fmt.Sprintf("%s-%s", info.ID, plan.ID)

	return &ServiceCatalogEntry{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Info:       info,
		Plan:       plan,
	}

}
