package mode

import (
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
