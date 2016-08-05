package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// SecretCreator creates a secret in the given namespace and returns either the created secret or a non-nil error
type SecretCreator interface {
	Create(string, *api.Secret) (*api.Secret, error)
}

// SecretDeleter deletes a secret in the given namespace of the given name and returns a non-nil error if the delete failed
type SecretDeleter interface {
	Delete(string, string) error
}

// SecretCreatorDeleter is a composition of SecretCreator and SecretDeleter. It's used for convenience to reduce 2 function parameters to 1
type SecretCreatorDeleter interface {
	SecretCreator
	SecretDeleter
}

type secCreatorDeleterImpl struct {
	creator SecretCreator
	deleter SecretDeleter
}

// NewSecretCreatorDeleter creates a new SecretCreatorDeleter from the given kcl.Client
func NewSecretCreatorDeleter(cl *kcl.Client) SecretCreatorDeleter {
	return secCreatorDeleterImpl{
		creator: NewSecretCreator(cl),
		deleter: NewSecretDeleter(cl),
	}
}

func (s secCreatorDeleterImpl) Create(ns string, sec *api.Secret) (*api.Secret, error) {
	return s.creator.Create(ns, sec)
}

func (s secCreatorDeleterImpl) Delete(ns, name string) error {
	return s.deleter.Delete(ns, name)
}

type secCreatorImpl struct {
	nsr kcl.SecretsNamespacer
}

func (s secCreatorImpl) Create(ns string, sec *api.Secret) (*api.Secret, error) {
	return s.nsr.Secrets(ns).Create(sec)
}

// NewSecretCreator is a convenience function to create a SecretCreator from a SecretsNamespacer
func NewSecretCreator(secretsNamespacer kcl.SecretsNamespacer) SecretCreator {
	return secCreatorImpl{nsr: secretsNamespacer}
}

type secDeleterImpl struct {
	nsr kcl.SecretsNamespacer
}

func (s secDeleterImpl) Delete(ns, name string) error {
	return s.nsr.Secrets(ns).Delete(name)
}

// NewSecretDeleter is a convenience function to create a SecretDeleter from a SecretsNamespacer
func NewSecretDeleter(secretsNamespacer kcl.SecretsNamespacer) SecretDeleter {
	return secDeleterImpl{nsr: secretsNamespacer}
}

// FakeSecretCreator is a fake implementation of SecretCreator. It can be used as a mock for tests
type FakeSecretCreator struct{}

// Create is the SecretCreator interface implementation. It copies s, sets its namespace to ns, and returns the copy and a non-nil error
func (f FakeSecretCreator) Create(ns string, s *api.Secret) (*api.Secret, error) {
	ret := *s
	s.ObjectMeta.Namespace = ns
	return &ret, nil
}

// FakeSecretDeleter is a fake implementation of SecretDeleter. It can be used as a mock for tests
type FakeSecretDeleter struct{}

// Delete is the SecretDeleter interface implementation. It simply returns nil
func (f FakeSecretDeleter) Delete(ns, name string) error {
	return nil
}
