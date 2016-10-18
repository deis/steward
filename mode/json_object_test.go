package mode

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/pborman/uuid"
)

func TestJSONObjectString(t *testing.T) {
	t.Skip("TODO")
}

func TestJSONObjectToFromStringRoundTrip(t *testing.T) {
	t.Run("full JSONObject", func(t *testing.T) {
		jso := JSONObject(map[string]interface{}{
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
		jso := JSONObject(map[string]interface{}{})
		jsoStr := jso.EncodeToString()
		jsoDecoded, err := JSONObjectFromString(jsoStr)
		assert.NoErr(t, err)
		assert.Equal(t, len(jsoDecoded), len(jso), "decoded json object length")
	})
}
