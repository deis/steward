package k8s

// ServiceCatalogLookup is an an O(1) lookup table for service catalog entries, based on ServiceCatalogEntry names. None of the functions on this struct are thread-safe
type ServiceCatalogLookup struct {
	lookup map[string]*ServiceCatalogEntry
}

// NewServiceCatalogLookup creates a new ServiceCatalogLookup, with its internal lookup table filled with all items in catalog.Items
func NewServiceCatalogLookup(catalog []*ServiceCatalogEntry) ServiceCatalogLookup {
	set := map[string]*ServiceCatalogEntry{}
	for _, item := range catalog {
		set[item.ResourceName()] = item
	}
	return ServiceCatalogLookup{lookup: set}
}

// Get looks up entry from the internal lookup table. Returns it if it exists, or nil if it doesn't
func (scl ServiceCatalogLookup) Get(entry *ServiceCatalogEntry) *ServiceCatalogEntry {
	ret, ok := scl.lookup[entry.ResourceName()]
	if !ok {
		return nil
	}
	return ret
}

// Set sets entry into the internal lookup table. Overwrites entries if they already existed
func (scl *ServiceCatalogLookup) Set(entry *ServiceCatalogEntry) {
	scl.lookup[entry.ResourceName()] = entry
}
