package k8s

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	serviceCatalogEntries3PRURLName = "servicecatalogentries"
	serviceCatalogEntryNamespace    = "steward"
)

type k8sRestClientImpl struct {
	cl *restclient.RESTClient
}

// NewK8sServiceCatalogInteractor creates a new ServiceCatalogInteractor which uses the Kubernetes API (using cl) to implement its functionality
func NewK8sServiceCatalogInteractor(cl *restclient.RESTClient) ServiceCatalogInteractor {
	return k8sRestClientImpl{cl: cl}
}

func (k k8sRestClientImpl) List() (*ServiceCatalogEntryList, error) {
	req := k.cl.Get().AbsPath(getServiceCatalogEntriesAbsPath()...)
	logger.Debugf("making request to %s", req.URL().String())
	res := req.Do()
	if res.Error() != nil {
		return nil, res.Error()
	}
	b, err := res.Raw()
	if err != nil {
		return nil, err
	}
	catalog := new(ServiceCatalogEntryList)
	if err := json.Unmarshal(b, catalog); err != nil {
		return nil, err
	}
	return catalog, nil
}

func (k k8sRestClientImpl) Create(entry *ServiceCatalogEntry) (*ServiceCatalogEntry, error) {
	b, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	res := k.cl.Post().AbsPath(getServiceCatalogEntriesAbsPath()...).Body(b).Do()
	if res.Error() != nil {
		return nil, res.Error()
	}
	return entry, nil
}

func getServiceCatalogEntriesAbsPath() []string {
	return []string{"apis", resourceAPIVersionBase, apiVersionV1, "namespaces", serviceCatalogEntryNamespace, serviceCatalogEntries3PRURLName}
}
