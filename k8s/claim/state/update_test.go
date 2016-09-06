package state

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
)

func TestUpdateClaim(t *testing.T) {
	type testCase struct {
		claim  mode.ServicePlanClaim
		update Update
	}
	testCases := []testCase{
		testCase{
			claim: mode.ServicePlanClaim{
				Status:            mode.StatusBinding.String(),
				StatusDescription: "some description",
			},
			update: NewUpdate(mode.StatusBound, "some other description", mode.JSONObject(map[string]string{"a": "b"})),
		},
		testCase{
			claim: mode.ServicePlanClaim{
				Status:            mode.StatusProvisioned.String(),
				StatusDescription: "start",
				Extra:             mode.JSONObject(map[string]string{"a": "b"}),
			},
			update: NewUpdate(mode.StatusBinding, "end", mode.JSONObject(map[string]string{"c": "d", "e": "f"})),
		},
	}

	for _, testCase := range testCases {
		UpdateClaim(&testCase.claim, testCase.update)
		assert.Equal(t, mode.Status(testCase.claim.Status), testCase.update.NewStatus, "new status")
		assert.Equal(t, testCase.claim.StatusDescription, testCase.update.NewStatusDescription, "new status description")
		assert.Equal(t, len(testCase.claim.Extra), len(testCase.update.NewExtra), "extra")
		for k, v := range testCase.claim.Extra {
			assert.Equal(t, testCase.update.NewExtra[k], v, "value of key "+k)
		}
	}
}
