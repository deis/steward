package brokerapi

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
)

const (
	ns   = "testnamespace"
	name = "testname"
)

func TestGetObjectMeta(t *testing.T) {
	meta := getObjectMeta(ns, name)
	assert.Equal(t, meta.Labels["created-by"], "steward", "created-by label")
	assert.Equal(t, meta.Namespace, ns, "namespace")
	assert.Equal(t, meta.Name, name, "resource name")
}

func TestWriteToKubernetes(t *testing.T) {
	creds := mode.JSONObject(map[string]string{
		"username": "testuser",
		"password": "testpass",
		"key":      "testkey",
	})
	cmCreator := &k8s.FakeConfigMapCreator{}
	assert.NoErr(t, writeToKubernetes(ns, name, creds, cmCreator))
	assert.Equal(t, len(cmCreator.Created), 1, "number of created ConfigMaps")
	cm := cmCreator.Created[0]
	assert.Equal(t, cm.Name, name, "ConfigMap name")
	assert.Equal(t, cm.Namespace, ns, "ConfigMap namespace")
	assert.Equal(t, cm.Labels["created-by"], "steward", "created-by label")
	assert.Equal(t, len(cm.Data), len(creds), "amount of stored data")
}
