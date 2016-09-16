package k8s

import (
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/watch"
)

// FakeConfigMapsInterface is a fake version of (k8s.io/client-go/1.4/kubernetes/typed/core/v1).ConfigMapInterface, for use in unit tests
type FakeConfigMapsInterface struct {
	Created    []*v1.ConfigMap
	GetReturns map[string]*v1.ConfigMap
}

// NewFakeConfigMapsInterface returns a new, empty *FakeConfigMapsInterface
func NewFakeConfigMapsInterface() *FakeConfigMapsInterface {
	return &FakeConfigMapsInterface{Created: nil, GetReturns: make(map[string]*v1.ConfigMap)}
}

// Get is the ConfigMapInterface interface implementation. If name is in f.GetReturns, returns f.GetReturns[name], nil. Otherwise returns nil, nil
func (f *FakeConfigMapsInterface) Get(name string) (*v1.ConfigMap, error) {
	cm, ok := f.GetReturns[name]
	if ok {
		return cm, nil
	}
	return nil, nil
}

// List is the ConfigMap interface implementation. It currently is not implemented and returns nil, nil
func (f *FakeConfigMapsInterface) List(opts api.ListOptions) (*v1.ConfigMapList, error) {
	return nil, nil
}

// Create is the ConfigMapInterface interface implementation. It appends cm to f.Created and then returns cm, nil. This function is not concurrency-safe
func (f *FakeConfigMapsInterface) Create(cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	f.Created = append(f.Created, cm)
	return cm, nil
}

// Delete is the ConfigMapInterface interface implementation. It currently is not implemented and returns nil.
func (f *FakeConfigMapsInterface) Delete(string, *api.DeleteOptions) error {
	return nil
}

// DeleteCollection is the ConfigMapInterface interface implementation. It current is not implemented and returns nil
func (f *FakeConfigMapsInterface) DeleteCollection(*api.DeleteOptions, api.ListOptions) error {
	return nil
}

// Update is the ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil.
func (f *FakeConfigMapsInterface) Update(*v1.ConfigMap) (*v1.ConfigMap, error) {
	return nil, nil
}

// Patch is the ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil.
func (f *FakeConfigMapsInterface) Patch(string, api.PatchType, []byte, ...string) (*v1.ConfigMap, error) {
	return nil, nil
}

// Watch is the ConfigMapsInterface interface implementation. It currently is not implemented and returns nil, nil.
func (f *FakeConfigMapsInterface) Watch(api.ListOptions) (watch.Interface, error) {
	return nil, nil
}
