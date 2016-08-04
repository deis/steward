package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// SecretCreator creates a secret in the given namespace and returns either the created secret or a non-nil error
type SecretCreator func(string, *api.Secret) (*api.Secret, error)

// NewSecretCreator is a convenience function to create a SecretCreator from a SecretsNamespacer
func NewSecretCreator(secretsNamespacer kcl.SecretsNamespacer) SecretCreator {
	return SecretCreator(func(ns string, sec *api.Secret) (*api.Secret, error) {
		return secretsNamespacer.Secrets(ns).Create(sec)
	})
}

// FakeSecretCreator returns a ConfigMapCreator that simply returns the given config map with its namespace set to the given namespace
func FakeSecretCreator() SecretCreator {
	return SecretCreator(func(ns string, s *api.Secret) (*api.Secret, error) {
		ret := *s
		s.ObjectMeta.Namespace = ns
		return &ret, nil
	})
}
