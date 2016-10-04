package k8s

import (
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/watch"
)

// FakeSecretsInterface is a fake version of (k8s.io/client-go/1.4/kubernetes/typed/core/v1).SecretInterface, for use in unit tests
type FakeSecretsInterface struct {
	Created []*v1.Secret
	Deleted []string
}

// Get is the SecretInterface interface implementation. It currently is not implemented and returns nil, nil
func (f *FakeSecretsInterface) Get(string) (*v1.Secret, error) {
	return nil, nil
}

// List is the SecretInterface interface implementation. It currently is not implemented and returns nil, nil
func (f *FakeSecretsInterface) List(opts api.ListOptions) (*v1.SecretList, error) {
	return nil, nil
}

// Create is the SecretInterface interface implementation. It appends cm to f.Created and then returns cm, nil. This function is not concurrency-safe
func (f *FakeSecretsInterface) Create(secret *v1.Secret) (*v1.Secret, error) {
	f.Created = append(f.Created, secret)
	return secret, nil
}

// Delete is the SecretInterface interface implementation. It appends name to f.Deleted and returns nil. This function is not concurrency-safe
func (f *FakeSecretsInterface) Delete(name string, opts *api.DeleteOptions) error {
	f.Deleted = append(f.Deleted, name)
	return nil
}

// DeleteCollection is the SecretInterface interface implementation. It currently is not implemented and returns nil
func (f *FakeSecretsInterface) DeleteCollection(*api.DeleteOptions, api.ListOptions) error {
	return nil
}

// Update is the SecretInterface interface implementation. It currently is not implemented and returns nil, nil
func (f *FakeSecretsInterface) Update(*v1.Secret) (*v1.Secret, error) {
	return nil, nil
}

// Patch is the SecretInterface interface implementation. It currently is not implemented and returns nil, ni
func (f FakeSecretsInterface) Patch(string, api.PatchType, []byte, ...string) (*v1.Secret, error) {
	return nil, nil
}

// Watch is the (k8s.io/kubernetes/pkg/client/unversioned).Secretsnterface interface implementation. It currently is not implemented and returns nil, nil
func (f *FakeSecretsInterface) Watch(api.ListOptions) (watch.Interface, error) {
	return nil, nil
}
