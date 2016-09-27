// +build integration

package cmd

import (
	"os"
	"testing"

	"github.com/deis/steward/test-utils/k8s"
	testsetup "github.com/deis/steward/test-utils/setup"
	"github.com/technosophos/moniker"
)

var namespace string

func TestMain(m *testing.M) {
	testsetup.SetupAndTearDown(m, setup, teardown)
}

func setup() error {
	namespace = moniker.New().NameSep("-")
	if err := k8s.EnsureNamespace(namespace); err != nil {
		return err
	}
	os.Setenv("POD_NAMESPACE", namespace)
	os.Setenv("CMD_IMAGE", "quay.io/deisci/cmd-sample-broker:devel")
	return nil
}

func teardown() error {
	if err := k8s.DeleteNamespace(namespace); err != nil {
		return err
	}
	return nil
}
