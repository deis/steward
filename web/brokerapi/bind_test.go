package brokerapi

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/arschles/testsrv"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
	"github.com/deis/steward/testutil"
	"github.com/pborman/uuid"
)

func TestBind(t *testing.T) {
	instID, bindID := uuid.New(), uuid.New()

	testutil.ConfigLogger()
	auth := testutil.GetAuth()
	bindCreds := map[string]string{
		"cred-1": "cred1",
		"cred-2": "cred2",
		"cred-3": "cred3",
	}
	backendBroker := testutil.NewCFBroker(auth, bindCreds)
	backendBrokerHost, backendBrokerPort, err := testutil.HostAndPort(backendBroker)
	assert.NoErr(t, err)

	cmNamespacer := k8s.NewFakeConfigMapsNamespacer()
	cfClient := cf.NewRESTClient(http.DefaultClient, "http", backendBrokerHost, backendBrokerPort, auth.Username, auth.Password)
	lifecycler := cf.NewLifecycler(cfClient)
	hdl := Handler(nil, lifecycler, auth, cmNamespacer)
	srv := testsrv.StartServer(hdl)
	defer srv.Close()

	reqBody := mode.BindRequest{
		ServiceID: uuid.New(),
		PlanID:    uuid.New(),
		Parameters: mode.JSONObject(map[string]string{
			"target_namespace": ns,
			"target_name":      name,
		}),
	}
	w := httptest.NewRecorder()
	req, err := testutil.NewReq(srv, auth, "PUT", nil, reqBody, "v2", "service_instances", instID, "service_bindings", bindID)
	assert.NoErr(t, err)
	hdl.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK, "response code")

	// check the k8s ConfigMapCreator's call
	assert.Equal(t, len(cmNamespacer.Returned), 1, "number of calls to the config map creator")
	cmInterface, ok := cmNamespacer.Returned[ns]
	assert.True(t, ok, "no config maps interface was returned for namespace %s", ns)
	assert.Equal(t, len(cmInterface.Created), 1, "number of config maps created")
	created := cmInterface.Created[0]
	assert.Equal(t, created.Namespace, ns, "config map namespace")
	assert.Equal(t, created.Name, name, "config map name")
	assert.Equal(t, len(created.Data), len(bindCreds), "amount of data in config map")
	for k, v := range bindCreds {
		assert.Equal(t, created.Data[k], base64.StdEncoding.EncodeToString([]byte(v)), "value of key "+k)
	}
}
