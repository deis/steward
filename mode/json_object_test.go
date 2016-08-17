package mode

import (
	"encoding/base64"
	"testing"

	"github.com/arschles/assert"
	"github.com/pborman/uuid"
)

func TestJSONObjectString(t *testing.T) {
	t.Skip("TODO")
}

func TestJSONObjectBase64EncodedVals(t *testing.T) {
	obj := JSONObject(map[string]string{
		"k1":       uuid.New(),
		uuid.New(): "v2",
	})

	encoded := obj.Base64EncodedVals()
	decodedMap := map[string]string{}
	for k, v := range encoded {
		decodedBytes, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			t.Fatalf("decoding failed for key %s, val %s (%s)", k, v, err)
		}
		decodedMap[k] = string(decodedBytes)
	}

	assert.Equal(t, len(decodedMap), len(obj), "length of JSONObjects")
	for k, v := range obj {
		assert.Equal(t, decodedMap[k], v, "value of key "+k)
	}
}
