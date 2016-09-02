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

func TestJSONObjectToFromStringRoundTrip(t *testing.T) {
	t.Run("full JSONObject", func(t *testing.T) {
		jso := JSONObject(map[string]string{
			"key1":     "val1",
			"key2":     "val2",
			"key3":     "val3",
			uuid.New(): uuid.New(),
		})
		jsoStr := jso.EncodeToString()
		jsoDecoded, err := JSONObjectFromString(jsoStr)
		assert.NoErr(t, err)
		assert.Equal(t, len(jsoDecoded), len(jso), "decoded json object length")
	})
	t.Run("empty JSONObject", func(t *testing.T) {
		jso := JSONObject(map[string]string{})
		jsoStr := jso.EncodeToString()
		jsoDecoded, err := JSONObjectFromString(jsoStr)
		assert.NoErr(t, err)
		assert.Equal(t, len(jsoDecoded), len(jso), "decoded json object length")
	})
}
