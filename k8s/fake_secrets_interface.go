package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

// FakeSecretsInterface is a fake version of (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface, for use in unit tests
type FakeSecretsInterface struct {
	Created []*api.Secret
}

// Get is the (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeSecretsInterface) Get(string) (*api.Secret, error) {
	return nil, nil
}

// List is the (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeSecretsInterface) List(opts api.ListOptions) (*api.SecretList, error) {
	return nil, nil
}

// Create is the (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface interface implementation. It appends cm to f.Created and then returns cm, nil. This function is not concurrency-safe
func (f *FakeSecretsInterface) Create(secret *api.Secret) (*api.Secret, error) {
	f.Created = append(f.Created, secret)
	return secret, nil
}

// Delete is the (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface interface implementation. It currently is not implemented and returns nil. It will be implemented if the need arises in tests
func (f *FakeSecretsInterface) Delete(string) error {
	return nil
}

// Update is the (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeSecretsInterface) Update(*api.Secret) (*api.Secret, error) {
	return nil, nil
}

// Watch is the (k8s.io/kubernetes/pkg/client/unversioned).Secretsnterface interface implementation. It currently is not implemented and returns nil, nil. It will be implemented if the need arises in tests
func (f *FakeSecretsInterface) Watch(api.ListOptions) (watch.Interface, error) {
	return nil, nil
}
