package k8s

import (
	"fmt"

	"k8s.io/kubernetes/pkg/client/restclient"
)

// ServiceCatalogLookup is an an O(1) lookup table for service catalog entries, based on ServiceCatalogEntry names. None of the functions on this struct are thread-safe
type ServiceCatalogLookup struct {
	lookup map[string]*ServiceCatalogEntry
}

// NewServiceCatalogLookup creates a new ServiceCatalogLookup, with its internal lookup table filled with all items in catalog.Items
func NewServiceCatalogLookup(catalog []*ServiceCatalogEntry) ServiceCatalogLookup {
	set := map[string]*ServiceCatalogEntry{}
	for _, item := range catalog {
		set[catalogKey(item.Info.ID, item.Plan.ID)] = item
	}
	return ServiceCatalogLookup{lookup: set}
}

// FetchServiceCatalogLookup returns a new ServiceCatalogLookup from the Kubernetes cluster using cl. Returns a non-nil error if there was a problem communicating with the cluster
func FetchServiceCatalogLookup(cl *restclient.RESTClient) (*ServiceCatalogLookup, error) {
	ret := NewServiceCatalogLookup(nil)
	entries, err := getServiceCatalogEntries(cl)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		ret.Set(entry)
	}
	return &ret, nil
}

// Get looks up entry from the internal lookup table. Returns it if it exists, or nil if it doesn't
func (scl ServiceCatalogLookup) Get(svcID, planID string) *ServiceCatalogEntry {
	ret, ok := scl.lookup[catalogKey(svcID, planID)]
	if !ok {
		return nil
	}
	return ret
}

// Set sets entry into the internal lookup table. Overwrites entries if they already existed
func (scl *ServiceCatalogLookup) Set(entry *ServiceCatalogEntry) {
	scl.lookup[catalogKey(entry.Info.ID, entry.Plan.ID)] = entry
}

func (scl *ServiceCatalogLookup) Len() int {
	return len(scl.lookup)
}

func catalogKey(svcID, planID string) string {
	return fmt.Sprintf("%s-%s", svcID, planID)
}
