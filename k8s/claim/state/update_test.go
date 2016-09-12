package state

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
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
			update: FullUpdate(
				mode.StatusBound,
				"some other description",
				uuid.New(),
				uuid.New(),
				mode.JSONObject(map[string]string{"a": "b"}),
			),
		},
		testCase{
			claim: mode.ServicePlanClaim{
				Status:            mode.StatusProvisioned.String(),
				StatusDescription: "start",
				Extra:             mode.JSONObject(map[string]string{"a": "b"}),
			},
			update: FullUpdate(
				mode.StatusBinding,
				"end",
				uuid.New(),
				uuid.New(),
				mode.JSONObject(map[string]string{"c": "d", "e": "f"}),
			),
		},
		testCase{
			claim: mode.ServicePlanClaim{
				Status:            mode.StatusProvisioned.String(),
				StatusDescription: "something",
				Extra:             mode.EmptyJSONObject(),
			},
			update: ErrUpdate(errors.New("error")),
		},
		testCase{
			claim: mode.ServicePlanClaim{
				Status:            mode.StatusProvisioned.String(),
				StatusDescription: "something else",
				Extra:             mode.EmptyJSONObject(),
			},
			update: StatusUpdate(mode.StatusBinding),
		},
	}

	for _, testCase := range testCases {
		origClaim := testCase.claim
		UpdateClaim(&testCase.claim, testCase.update)
		_, isFullUpdate := testCase.update.(fullUpdate)
		eUpdate, isErrUpdate := testCase.update.(errUpdate)
		_, isStatusUpdate := testCase.update.(statusUpdate)

		assert.Equal(t, mode.Status(testCase.claim.Status), testCase.update.Status(), "new status")
		if isFullUpdate {
			// on full update, the description should be updated
			assert.Equal(t, testCase.claim.StatusDescription, testCase.update.Description(), "new status description")
		} else if isErrUpdate {
			// on error update, the description should be the error string
			assert.Equal(t, testCase.claim.StatusDescription, eUpdate.err.Error(), "new status description")
		} else if isStatusUpdate {
			// on status update, the description should be unchanged
			assert.Equal(t, testCase.claim.StatusDescription, origClaim.StatusDescription, "new status description")
		}
		if isFullUpdate {
			// on full update, the extra should be set
			assert.Equal(t, len(testCase.claim.Extra), len(testCase.update.Extra()), "extra")
			newExtra := testCase.update.Extra()
			for k, v := range testCase.claim.Extra {
				assert.Equal(t, newExtra[k], v, "value of key "+k)
			}
		} else {
			// otherwise, the extra should be unchanged
			assert.Equal(t, len(testCase.claim.Extra), len(origClaim.Extra), "extra")
			oldExtra := origClaim.Extra
			for k, v := range testCase.claim.Extra {
				assert.Equal(t, oldExtra[k], v, "value of key "+k)
			}
		}
	}
}
