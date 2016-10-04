package k8s

import (
	"sync"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/errors"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/rest"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

const (
	kubeConfigFile = "/go/src/github.com/deis/steward/kubeconfig.yaml"
)

var (
	clientset  *kubernetes.Clientset
	clientOnce sync.Once
)

// GetClientset returns a *kubernetes.Clientset suitable for use with integration tests against
// a leased k8s cluster
func GetClientset() (*kubernetes.Clientset, error) {
	var err error
	clientOnce.Do(func() {
		var config *rest.Config
		if config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile); err == nil {
			clientset, err = kubernetes.NewForConfig(config)
		}
	})
	return clientset, err
}

// EnsureNamespace ensures the existence of the specified namespace
func EnsureNamespace(namespaceStr string) error {
	clientset, err := GetClientset()
	if err != nil {
		return err
	}
	nsClient := clientset.Namespaces()
	// Just try to create the namespace. If it already exists, that's fine.
	_, err = nsClient.Create(&v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: namespaceStr,
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// DeleteNamespace deletes the specified namespace
func DeleteNamespace(namespaceStr string) error {
	clientset, err := GetClientset()
	if err != nil {
		return err
	}
	nsClient := clientset.Namespaces()
	// If the problem is just that the namespace doesn't exist, ignore it
	if err := nsClient.Delete(namespaceStr, &api.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
