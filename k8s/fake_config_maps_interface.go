package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

// FakeConfigMapsInterface is a fake version of (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface, for use in unit tests
type FakeConfigMapsInterface struct {
	Created []*api.ConfigMap
}

// Get is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeConfigMapsInterface) Get(string) (*api.ConfigMap, error) {
	return nil, nil
}

// List is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeConfigMapsInterface) List(opts api.ListOptions) (*api.ConfigMapList, error) {
	return nil, nil
}

// Create is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface interface implementation. It appends cm to f.Created and then returns cm, nil. This function is not concurrency-safe
func (f *FakeConfigMapsInterface) Create(cm *api.ConfigMap) (*api.ConfigMap, error) {
	f.Created = append(f.Created, cm)
	return cm, nil
}

// Delete is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface interface implementation. It currently is not implemented and returns nil. It will be implemented if the need arises in tests
func (f *FakeConfigMapsInterface) Delete(string) error {
	return nil
}

// Update is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeConfigMapsInterface) Update(*api.ConfigMap) (*api.ConfigMap, error) {
	return nil, nil
}

// Watch is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeConfigMapsInterface) Watch(api.ListOptions) (watch.Interface, error) {
	return nil, nil
}
