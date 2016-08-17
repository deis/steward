package k8s

import (
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// FakeConfigMapsNamespacer is a fake implementation of (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsNamespacer, suitable for use in unit tests.
type FakeConfigMapsNamespacer struct {
	Returned map[string]*FakeConfigMapsInterface
}

// NewFakeConfigMapsNamespacer returns an empty FakeConfigMapsNamespacer
func NewFakeConfigMapsNamespacer() *FakeConfigMapsNamespacer {
	return &FakeConfigMapsNamespacer{Returned: make(map[string]*FakeConfigMapsInterface)}
}

// ConfigMaps is the (k8s.io/kubernetes/pkg/client/unversioned).ConfigMapsNamespacer interface implementation. It returns an empty kcl.ConfigMapsInterface
func (f *FakeConfigMapsNamespacer) ConfigMaps(ns string) kcl.ConfigMapsInterface {
	ret := &FakeConfigMapsInterface{}
	f.Returned[ns] = ret
	return ret
}
