package k8s

import (
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
)

type jsonObjectEqualer mode.JSONObject

func (j jsonObjectEqualer) Equal(e assert.Equaler) bool {
	je, ok := e.(jsonObjectEqualer)
	if !ok {
		return false
	}
	thisObj := mode.JSONObject(j)
	otherObj := mode.JSONObject(je)
	if len(thisObj) != len(otherObj) {
		return false
	}

	for key, val := range thisObj {
		if val != otherObj[key] {
			return false
		}
	}
	return true
}

type servicePlanClaimEqualer ServicePlanClaim

func (s servicePlanClaimEqualer) Equal(e assert.Equaler) bool {
	thisObj := ServicePlanClaim(s)
	spce, ok := e.(servicePlanClaimEqualer)
	if !ok {
		return false
	}
	otherObj := ServicePlanClaim(spce)
	if thisObj.TargetName != otherObj.TargetName ||
		thisObj.ServiceID != otherObj.ServiceID ||
		thisObj.PlanID != otherObj.PlanID ||
		thisObj.ClaimID != otherObj.ClaimID ||
		thisObj.Action != otherObj.Action ||
		thisObj.Status != otherObj.Status ||
		thisObj.StatusDescription != otherObj.StatusDescription ||
		thisObj.InstanceID != otherObj.InstanceID ||
		thisObj.BindID != otherObj.BindID ||
		!jsonObjectEqualer(thisObj.Extra).Equal(jsonObjectEqualer(otherObj.Extra)) {
		return false
	}
	return true
}

func TestErrDataMapMissingKey(t *testing.T) {
	err := errDataMapMissingKey{key: "testKey"}
	assert.True(t, strings.Contains(err.Error(), err.key), "error string didn't contain error key %s", err.key)
}

func TestServicePlanClaimMapRoundTrip(t *testing.T) {
	claim := ServicePlanClaim{
		TargetName:        "testTarget",
		ServiceID:         uuid.New(),
		PlanID:            uuid.New(),
		ClaimID:           uuid.New(),
		Action:            "testAction",
		Status:            "testStatus",
		StatusDescription: "testStatusDescription",
		InstanceID:        uuid.New(),
		BindID:            uuid.New(),
		Extra:             mode.JSONObject(map[string]string{"key1": "val1"}),
	}
	m := claim.ToMap()
	parsedClaim, err := ServicePlanClaimFromMap(m)
	assert.NoErr(t, err)
	// note that servicePlanClaimEqualer is an instance of assert.Equals, and assert.Equal checks for instances of that type (https://godoc.org/github.com/arschles/assert#Equaler). This feature allows us to define what "deep equals" means (instead of using reflect.DeepEqual)
	assert.Equal(t, servicePlanClaimEqualer(*parsedClaim), servicePlanClaimEqualer(claim), "parsed claim")
}
