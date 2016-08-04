package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// ConfigMapCreator creates a config map in the given namespace and returns either the created config map or a non-nil error
type ConfigMapCreator func(string, *api.ConfigMap) (*api.ConfigMap, error)

// NewConfigMapCreator is a convenience function to create a ConfigMapCreator from a ConfigMapsNamespacer
func NewConfigMapCreator(configMapsNamespacer kcl.ConfigMapsNamespacer) ConfigMapCreator {
	return ConfigMapCreator(func(ns string, cm *api.ConfigMap) (*api.ConfigMap, error) {
		return configMapsNamespacer.ConfigMaps(ns).Create(cm)
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
