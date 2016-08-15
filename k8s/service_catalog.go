package k8s

import (
	"encoding/json"
	"fmt"

	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	serviceCatalogEntry3PRName      = "ServiceCatalogEntry"
	serviceCatalogEntries3PRURLName = "servicecatalogentries"
	serviceCatalogEntryNamespace    = "steward"
)

// ServiceCatalogEntry is the third party resource data that represents a single service provider + plan. note that this structure is not written directly into the third party.
type ServiceCatalogEntry struct {
	Info mode.ServiceInfo `json:"service_info"`
	Plan mode.ServicePlan `json:"service_plan"`
}

// TODO: make this conform to runtime.Object
type serviceCatalogEntry3PRWrapper struct {
	api.ObjectMeta       `json:"metadata,omitempty"`
	unversioned.TypeMeta `json:",inline"`
	*ServiceCatalogEntry `json:",inline"`
}

// ResourceName returns the name of the ServiceProviderPlanPair third party resource
func (sce ServiceCatalogEntry) ResourceName() string {
	return fmt.Sprintf("%s-%s", sce.Info.Name, sce.Plan.Name)
}

type serviceCatalogEntries3PRWrapper struct {
	api.ObjectMeta       `json:"metadata,omitempty"`
	unversioned.TypeMeta `json:",inline"`
	Items                []*ServiceCatalogEntry `json:"items"`
}

// PublishServiceCatalogEntry publishes spp to the service catalog third party resource
func PublishServiceCatalogEntry(cl *restclient.RESTClient, spp *ServiceCatalogEntry) error {
	wrapper := serviceCatalogEntry3PRWrapper{
		ObjectMeta:          api.ObjectMeta{Name: spp.ResourceName()},
		TypeMeta:            unversioned.TypeMeta{Kind: serviceCatalogEntry3PRName, APIVersion: resourceAPIVersion(apiVersionV1)},
		ServiceCatalogEntry: spp,
	}
	// TODO: once serviceCatalogEntry3PRWrapper implements runtime.Object, remove this marshal and send wrapper directly
	b, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	res := cl.Post().AbsPath(getServiceCatalogEntriesAbsPath()...).Body(b).Do()
	if res.Error() != nil {
		return res.Error()
	}
	return nil
}

// GetServiceCatalogEntries gets a list of all services
func GetServiceCatalogEntries(
	cl *restclient.RESTClient,
) ([]*ServiceCatalogEntry, error) {

	req := cl.Get().AbsPath(getServiceCatalogEntriesAbsPath()...)
	logger.Debugf("making request to %s", req.URL().String())
	res := req.Do()
	if res.Error() != nil {
		return nil, res.Error()
	}
	b, err := res.Raw()
	if err != nil {
		return nil, err
	}
	catalog := new(serviceCatalogEntries3PRWrapper)
	if err := json.Unmarshal(b, catalog); err != nil {
		return nil, err
	}
	return catalog.Items, nil
}

func getServiceCatalogEntriesAbsPath() []string {
	return []string{"apis", "steward.deis.com", "v1", "namespaces", serviceCatalogEntryNamespace, serviceCatalogEntries3PRURLName}
}
