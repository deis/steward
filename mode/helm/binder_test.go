package helm

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/pborman/uuid"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/helm/pkg/chartutil"
)

func TestDataFieldKey(t *testing.T) {
	const (
		key = "key1"
	)
	cm := &v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: uuid.New(), Namespace: uuid.New()},
	}
	fullKey := dataFieldKey(cm, key)
	assert.Equal(t, fullKey, fmt.Sprintf("%s-%s-%s", cm.Namespace, cm.Name, key), "data field key")
}

func TestBind(t *testing.T) {
	const (
		instID      = "testInstID"
		bindID      = "testBindID"
		cmNamespace = "default"  // this is the creds config map namespace hard-coded in the alpine chart
		cmName      = "my-creds" // this is the creds config map name hard-coded in the alpine chart
	)
	chart, err := chartutil.Load(alpineChartLoc())
	assert.NoErr(t, err)
	nsr := k8s.NewFakeConfigMapsNamespacer()
	defaultIface := k8s.NewFakeConfigMapsInterface()
	defaultIface.GetReturns[cmName] = &v1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      "testName",
			Namespace: "testNamespace",
		},
		Data: map[string]string{"testKey": "testVal"},
	}
	nsr.ToReturn[cmNamespace] = defaultIface
	binder, err := newBinder(chart, nsr)
	assert.NoErr(t, err)
	resp, err := binder.Bind(instID, bindID, nil)
	assert.NoErr(t, err)
	assert.NotNil(t, resp, "bind response")
	assert.Equal(t, len(resp.Creds), len(defaultIface.GetReturns[cmName].Data), "length of returned credentials")
	expectedCM := defaultIface.GetReturns[cmName]
	for k, v := range expectedCM.Data {
		expectedKey := dataFieldKey(expectedCM, k)
		assert.Equal(t, v, resp.Creds[expectedKey], fmt.Sprintf("value of key %s", k))
	}
}
