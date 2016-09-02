package mode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	errMissing = errors.New("key is missing")
)

type errMalformedKV struct {
	kv []string
}

func (e errMalformedKV) Error() string {
	return fmt.Sprintf("malformed JSONObject-encoded key/value pair %s", e.kv)
}

// JSONObject is a convenience wrapper around a Go type that represents a JSON object
type JSONObject map[string]string

func (j JSONObject) String(key string) (string, error) {
	i, ok := j[key]
	if !ok {
		return "", errMissing
	}
	return i, nil
}

// Base64EncodedVals returns a new JSONObject equivalent to j with all values base64 Encoded
func (j JSONObject) Base64EncodedVals() JSONObject {
	newMap := make(map[string]string)
	for k, v := range j {
		newMap[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}
	return JSONObject(newMap)
}

// MarshalText is the encoding.TextMarshaler implementation
func (j JSONObject) EncodeToString() string {
	slc := make([]string, len(j))
	i := 0
	for key, val := range j {
		slc[i] = fmt.Sprintf("%s=%s", key, val)
		i++
	}
	return strings.Join(slc, ",")
}

// JSONObjectFromString decodes a string into a JSONObject. Returns a non-nil error if the string was not a valid JSONObject
func JSONObjectFromString(str string) (JSONObject, error) {
	if len(str) == 0 {
		return JSONObject(map[string]string{}), nil
	}
	mp := map[string]string{}
	spl := strings.Split(str, ",")
	for _, s := range spl {
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			return JSONObject(map[string]string{}), errMalformedKV{kv: kv}
		}
		key, val := kv[0], kv[1]
		mp[key] = val
	}
	return JSONObject(mp), nil
}
