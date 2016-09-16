package k8s

import (
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
)

// FakeConfigMapsNamespacer is a fake implementation of (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsNamespacer, suitable for use in unit tests.
type FakeConfigMapsNamespacer struct {
	ToReturn map[string]*FakeConfigMapsInterface
	Returned map[string]*FakeConfigMapsInterface
}

// NewFakeConfigMapsNamespacer returns an empty FakeConfigMapsNamespacer
func NewFakeConfigMapsNamespacer() *FakeConfigMapsNamespacer {
	return &FakeConfigMapsNamespacer{
		ToReturn: make(map[string]*FakeConfigMapsInterface),
		Returned: make(map[string]*FakeConfigMapsInterface),
	}
}

// ConfigMaps is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsNamespacer interface implementation. It returns an empty kcl.ConfigMapsInterface
func (f *FakeConfigMapsNamespacer) ConfigMaps(ns string) v1.ConfigMapInterface {
	iface, ok := f.ToReturn[ns]
	if !ok {
		iface = &FakeConfigMapsInterface{}
	}
	f.Returned[ns] = iface
	return iface
}
