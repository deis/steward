package k8s

// FakeServiceCatalogInteractor is a fake ServiceCatalogInteractor implementation, for use in unit tests
type FakeServiceCatalogInteractor struct {
	ListRet   *ServiceCatalogEntryList
	Created   []*ServiceCatalogEntry
	CreateErr error
}

// List is the ServiceCatalogInteractor interface implementation. It simply returns f.ListRet, nil
func (f FakeServiceCatalogInteractor) List() (*ServiceCatalogEntryList, error) {
	return f.ListRet, nil
}

// Create is the ServiceCatalogInteractor interface implementation. It appends entry to f.Created and returns entry, nil if f.CreateErr is nil. Otherwise, returns nil, f.CreateErr
func (f *FakeServiceCatalogInteractor) Create(entry *ServiceCatalogEntry) (*ServiceCatalogEntry, error) {
	f.Created = append(f.Created, entry)
	if f.CreateErr != nil {
		return nil, f.CreateErr
	}
	return entry, nil
}
