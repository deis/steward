package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// ConfigMapCreator creates a config map in the given namespace and returns either the created config map or a non-nil error
type ConfigMapCreator func(string, *api.ConfigMap) (*api.ConfigMap, error)

// ConfigMapDeleter deletes a config map from the given namespace. Returns a non-nil error if the delete failed
type ConfigMapDeleter func(string, string) error

// NewConfigMapCreator is a convenience function to create a ConfigMapCreator from a ConfigMapsNamespacer
func NewConfigMapCreator(configMapsNamespacer kcl.ConfigMapsNamespacer) ConfigMapCreator {
	return ConfigMapCreator(func(ns string, cm *api.ConfigMap) (*api.ConfigMap, error) {
		return configMapsNamespacer.ConfigMaps(ns).Create(cm)
	})
}

// NewConfigMapDeleter is a convenience function to create a ConfigMapDeleter from a ConfigMapsNamespacer
func NewConfigMapDeleter(configMapsNamespacer kcl.ConfigMapsNamespacer) ConfigMapDeleter {
	return ConfigMapDeleter(func(ns, name string) error {
		return configMapsNamespacer.ConfigMaps(ns).Delete(name)
	})
}

// FakeConfigMapCreator returns a ConfigMapCreator that simply returns the given config map with its namespace set to the given namespace
func FakeConfigMapCreator() ConfigMapCreator {
	return ConfigMapCreator(func(ns string, cm *api.ConfigMap) (*api.ConfigMap, error) {
		ret := *cm
		ret.ObjectMeta.Namespace = ns
		return &ret, nil
	})
}

// FakeConfigMapDeleter returns a ConfigMapDelete that does nothing and returns nil
func FakeConfigMapDeleter() ConfigMapDeleter {
	return ConfigMapDeleter(func(ns, name string) error {
		return nil
	})
}
