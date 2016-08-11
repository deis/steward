package brokerapi

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
)

const (
	ns     = "testnamespace"
	svcID  = "testsvc"
	planID = "testplan"
	bindID = "testbind"
	instID = "testinst"
)

func TestGetResourceName(t *testing.T) {
	name := getResourceName(svcID, planID, bindID, instID)
	assert.Equal(t, name, fmt.Sprintf("%s-%s-%s-%s", svcID, planID, bindID, instID), "resource name")
}

func TestGetObjectMeta(t *testing.T) {
	meta := getObjectMeta(ns, svcID, planID, bindID, instID)
	assert.Equal(t, meta.Labels["created-by"], "steward", "created-by label")
	assert.Equal(t, meta.Namespace, ns, "namespace")
	assert.Equal(t, meta.Name, getResourceName(svcID, planID, bindID, instID), "resource name")
}
