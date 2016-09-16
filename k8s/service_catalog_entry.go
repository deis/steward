package k8s

import (
	"fmt"
	"strings"

	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/unversioned"
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
	Description          string           `json:"description"`
}

// NewServiceCatalogEntry creates a new ServiceCatalogEntry suitable for writing to the Kubernetes API
func NewServiceCatalogEntry(
	brokerName string,
	objectMeta api.ObjectMeta,
	info mode.ServiceInfo,
	plan mode.ServicePlan) *ServiceCatalogEntry {
	typeMeta := unversioned.TypeMeta{
		Kind:       ServiceCatalogEntryKind,
		APIVersion: resourceAPIVersion(apiVersionV1),
	}
	canBrokerName := canonicalize(brokerName)
	canServiceName := canonicalize(info.Name)
	canPlanName := canonicalize(plan.Name)
	description := fmt.Sprintf("%s (%s)", info.Description, plan.Description)

	objectMeta.Name = fmt.Sprintf("%s-%s-%s", canBrokerName, canServiceName, canPlanName)

	objectMeta.Labels = map[string]string{
		"broker":       brokerName,
		"service-id":   info.ID,
		"service-name": info.Name,
		"plan-id":      plan.ID,
		"plan-name":    plan.Name,
	}

	return &ServiceCatalogEntry{
		TypeMeta:    typeMeta,
		ObjectMeta:  objectMeta,
		Info:        info,
		Plan:        plan,
		Description: description,
	}

}

// canonicalize transforms a given name into a name that doesn't contain characters that k8s
// doesn't permit in a resource name.
func canonicalize(name string) string {
	// krancour: This approach seemed really naive at first. It seemed there must WAY more things we
	// need to account for. However, after reading how the CF broker API's documenation describes
	// the service name and plan name fields, I am convinced that this is a reasonable start, which
	// can be amended later if necessary.
	//
	// The description:
	//
	// > The CLI-friendly name of the service that will appear in the catalog. All lowercase, no
	// > spaces.
	//
	// Although this doesn't explicitly forbid other characters, simplicity is strongly implied.
	// From experience, I know some brokers DO include other characters in service names or plan
	// names. Everything that follows is a concession to attempt to account for that.
	name = strings.Replace(name, "_", "-", -1)
	name = strings.Replace(name, ".", "-", -1)
	name = strings.Replace(name, ":", "-", -1)
	name = strings.Replace(name, "/", "-", -1)
	name = strings.Replace(name, "\\", "-", -1)
	// Per the comment above, spaces are not permitted, however not all brokers are CF brokers.
	// Although we'd like brokers to play by CF's rules, there is no guarantee.
	name = strings.Replace(name, " ", "-", -1)
	return name
}
