// +build integration

package cmd

import (
	"os"
	"testing"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/test-utils/k8s"
	testsetup "github.com/deis/steward/test-utils/setup"
	"github.com/technosophos/moniker"
	"k8s.io/client-go/1.4/kubernetes"
)

var (
	clientset      *kubernetes.Clientset
	testCataloger  mode.Cataloger
	testLifecycler *mode.Lifecycler
	testNamespace  string
)

func TestMain(m *testing.M) {
	testsetup.SetupAndTearDown(m, setup, teardown)
}

func setup() error {
	var err error
	if clientset, err = k8s.GetClientset(); err != nil {
		return err
	}
	testNamespace = moniker.New().NameSep("-")
	if err := k8s.EnsureNamespace(testNamespace); err != nil {
		return err
	}
	os.Setenv("POD_NAMESPACE", testNamespace)
	os.Setenv("CMD_IMAGE", "quay.io/deisci/cmd-sample-broker:devel")
	if testCataloger, testLifecycler, err = GetComponents(clientset); err != nil {
		return err
	}
	return nil
}

func teardown() error {
	if err := k8s.DeleteNamespace(testNamespace); err != nil {
		return err
	}
	return nil
}
