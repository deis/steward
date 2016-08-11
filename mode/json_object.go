package mode

import (
	"encoding/base64"
	"errors"
)

var (
	errMissing = errors.New("key is missing")
)

// JSONObject is a convenience wrapper around a Go type that represents a JSON object
type JSONObject map[string]string

func (j JSONObject) String(key string) (string, error) {
	i, ok := j[key]
	if !ok {
		return "", errMissing
	}
	return i, nil
}

// Exists returns whether key exists
func (j JSONObject) Exists(key string) bool {
	_, ok := j[key]
	return ok
}

// Base64EncodedVals returns a new JSONObject equivalent to j with all values base64 Encoded
func (j JSONObject) Base64EncodedVals() JSONObject {
	newMap := make(map[string]string)
	for k, v := range j {
		newMap[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}
	return JSONObject(newMap)
}
