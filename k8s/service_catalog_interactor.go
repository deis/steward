package k8s

// ServiceCatalogInteractor is a common set of functions for interacting with ServiceCatalogEntry third party resources in Kubernetes.
type ServiceCatalogInteractor interface {
	List() (*ServiceCatalogEntryList, error)
	Create(*ServiceCatalogEntry) (*ServiceCatalogEntry, error)
}
