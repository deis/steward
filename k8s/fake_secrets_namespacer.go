package k8s

import (
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// FakeSecretsNamespacer is a fake implementation of (k8s.io/kubernetes/pkg/client/unversioned).SecretsNamespacer, suitable for use in unit tests.
type FakeSecretsNamespacer struct {
	Returned map[string]*FakeSecretsInterface
}

// NewFakeSecretsNamespacer returns an empty FakeSecretsNamespacer
func NewFakeSecretsNamespacer() *FakeSecretsNamespacer {
	return &FakeSecretsNamespacer{Returned: make(map[string]*FakeSecretsInterface)}
}

// Secrets is the (k8s.io/kubernetes/pkg/client/unversioned).SecretsNamespacer interface implementation. It returns an empty kcl.SecretsInterface
func (f *FakeSecretsNamespacer) Secrets(ns string) kcl.SecretsInterface {
	ret := &FakeSecretsInterface{}
	f.Returned[ns] = ret
	return ret
}
