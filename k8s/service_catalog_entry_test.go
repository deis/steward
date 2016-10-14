package k8s

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"k8s.io/client-go/1.4/pkg/api"
)

func TestNewServiceCatalogEntry(t *testing.T) {
	brokerName := "testBroker"
	objectMeta := api.ObjectMeta{Name: "testOM", Namespace: "testNS"}
	info := mode.ServiceInfo{Name: "testServiceInfo", ID: uuid.New(), Description: "testDescrInfo"}
	plan := mode.ServicePlan{Name: "testPlan", ID: uuid.New(), Description: "testDescrPlan"}

	entry := NewServiceCatalogEntry(brokerName, objectMeta, info, plan)
	assert.Equal(t, entry.TypeMeta.Kind, ServiceCatalogEntryKind, "kind")
	assert.Equal(t, entry.TypeMeta.APIVersion, resourceAPIVersion(apiVersionV1), "API version")
	assert.Equal(t, entry.ObjectMeta.Name, fmt.Sprintf("%s-%s-%s", brokerName, info.Name, plan.Name), "object meta name")
	assert.Equal(t, entry.ObjectMeta.Labels["broker"], brokerName, "broker name label")
	assert.Equal(t, entry.ObjectMeta.Labels["service-id"], info.ID, "service ID label")
	assert.Equal(t, entry.ObjectMeta.Labels["plan-id"], plan.ID, "plan ID label")
	assert.Equal(t, entry.ObjectMeta.Labels["plan-name"], plan.Name, "plan name label")
	assert.Equal(t, entry.Info, info, "service info")
	assert.Equal(t, entry.Plan, plan, "service plan")
	assert.Equal(t, entry.Description, fmt.Sprintf("%s (%s)", info.Description, plan.Description), "service description")
}

func TestCanonicalize(t *testing.T) {
	strs := []string{
		"a_b_c",
		"a.b.c",
		"a:b:c",
		"a/b/c",
		`a\b\c`,
		"a b c",
	}
	for _, str := range strs {
		cleaned := canonicalize(str)
		assert.Equal(t, cleaned, "a-b-c", "canonicalized string")
	}
}
