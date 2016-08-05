package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// SecretCreator creates a secret in the given namespace and returns either the created secret or a non-nil error
type SecretCreator func(string, *api.Secret) (*api.Secret, error)

// SecretDeleter deletes a secret in the given namespace of the given name and returns a non-nil error if the delete failed
type SecretDeleter func(string, string) error

// NewSecretCreator is a convenience function to create a SecretCreator from a SecretsNamespacer
func NewSecretCreator(secretsNamespacer kcl.SecretsNamespacer) SecretCreator {
	return SecretCreator(func(ns string, sec *api.Secret) (*api.Secret, error) {
		return secretsNamespacer.Secrets(ns).Create(sec)
	})
}

// NewSecretDeleter is a convenience function to create a SecretDeleter from a SecretsNamespacer
func NewSecretDeleter(secretsNamespacer kcl.SecretsNamespacer) SecretDeleter {
	return SecretDeleter(func(ns, name string) error {
		return secretsNamespacer.Secrets(ns).Delete(name)
	})
}

// FakeSecretCreator returns a SecretCreator that simply returns the given secret with its namespace set to the given namespace
func FakeSecretCreator() SecretCreator {
	return SecretCreator(func(ns string, s *api.Secret) (*api.Secret, error) {
		ret := *s
		s.ObjectMeta.Namespace = ns
		return &ret, nil
	})
}

// FakeSecretDeleter returns a SecretDeleter does nothing and returns nil
func FakeSecretDeleter() SecretDeleter {
	return SecretDeleter(func(ns, name string) error {
		return nil
	})
}
