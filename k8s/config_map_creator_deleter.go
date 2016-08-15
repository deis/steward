package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// ConfigMapCreator creates a config map in the given namespace and returns either the created config map or a non-nil error
type ConfigMapCreator interface {
	Create(string, *api.ConfigMap) (*api.ConfigMap, error)
}

// ConfigMapDeleter deletes a config map from the given namespace. Returns a non-nil error if the delete failed
type ConfigMapDeleter interface {
	Delete(string, string) error
}

// ConfigMapCreatorDeleter is a composition of ConfigMapCreator and ConfigMapDeleter. It's used for convenience to reduce 2 function parameters to 1
type ConfigMapCreatorDeleter interface {
	ConfigMapCreator
	ConfigMapDeleter
}

type cmCreatorDeleterImpl struct {
	creator ConfigMapCreator
	deleter ConfigMapDeleter
}

// NewConfigMapCreatorDeleter creates a new SecretCreatorDeleter from the given kcl.Client
func NewConfigMapCreatorDeleter(cl *kcl.Client) ConfigMapCreatorDeleter {
	return cmCreatorDeleterImpl{
		creator: NewConfigMapCreator(cl),
		deleter: NewConfigMapDeleter(cl),
	}
}

func (s cmCreatorDeleterImpl) Create(ns string, sec *api.ConfigMap) (*api.ConfigMap, error) {
	return s.creator.Create(ns, sec)
}

func (s cmCreatorDeleterImpl) Delete(ns, name string) error {
	return s.deleter.Delete(ns, name)
}

type cmCreatorImpl struct {
	nsr kcl.ConfigMapsNamespacer
}

func (c cmCreatorImpl) Create(ns string, cm *api.ConfigMap) (*api.ConfigMap, error) {
	return c.nsr.ConfigMaps(ns).Create(cm)
}

// NewConfigMapCreator is a convenience function to create a ConfigMapCreator from a ConfigMapsNamespacer
func NewConfigMapCreator(configMapsNamespacer kcl.ConfigMapsNamespacer) ConfigMapCreator {
	return cmCreatorImpl{nsr: configMapsNamespacer}
}

type cmDeleterImpl struct {
	nsr kcl.ConfigMapsNamespacer
}

func (c cmDeleterImpl) Delete(ns, name string) error {
	return c.nsr.ConfigMaps(ns).Delete(name)
}

// NewConfigMapDeleter is a convenience function to create a ConfigMapDeleter from a ConfigMapsNamespacer
func NewConfigMapDeleter(configMapsNamespacer kcl.ConfigMapsNamespacer) ConfigMapDeleter {
	return cmDeleterImpl{nsr: configMapsNamespacer}
}

// FakeConfigMapCreator is a fake ConfigMapCreator, suitable for mocking in tests. It is not concurrency safe
type FakeConfigMapCreator struct {
	Created []*api.ConfigMap
}

// Create is the ConfigMapCreator interface implementation. It simply copies cm, changes its namespace to ns, and returns it along with a nil error
func (f *FakeConfigMapCreator) Create(ns string, cm *api.ConfigMap) (*api.ConfigMap, error) {
	ret := *cm
	ret.ObjectMeta.Namespace = ns
	f.Created = append(f.Created, &ret)
	return &ret, nil
}

// FakeConfigMapDeleter is a fake ConfigMapDeleter, suitable for mocking in tests
type FakeConfigMapDeleter struct{}

// Delete is the ConfigMapDeleter interface implementation. It simply returns nil
func (f *FakeConfigMapDeleter) Delete(ns, name string) error {
	return nil
}

// FakeConfigMapCreatorDeleter is a fake implementation of ConfigMapCreatorDeleter, suitable for mocking in tests
type FakeConfigMapCreatorDeleter struct {
	*FakeConfigMapCreator
	*FakeConfigMapDeleter
}
