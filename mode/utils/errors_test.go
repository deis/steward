package utils

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

func checkErr(e error, subStr string) error {
	if !strings.Contains(e.Error(), subStr) {
		return fmt.Errorf("error string doesn't contain substring '%s'", subStr)
	}
	return nil
}

func TestErrUnrecognizedMode(t *testing.T) {
	assert.NoErr(t, checkErr(errUnrecognizedMode{mode: "testmode"}, "testmode"))
}

func TestErrGettingK8sClient(t *testing.T) {
	assert.NoErr(t, checkErr(errGettingK8sClient{Original: errors.New("testorig")}, "testorig"))
}

func TestErrPublishingServiceCatalog(t *testing.T) {
	assert.NoErr(t, checkErr(errPublishingServiceCatalog{Original: errors.New("testorig")}, "testorig"))
}

func TestErrGettingServiceCatalogLookupTable(t *testing.T) {
	assert.NoErr(
		t,
		checkErr(errGettingServiceCatalogLookupTable{Original: errors.New("testorig")}, "testorig"),
	)
}

func TestErrGettingServiceCatalog(t *testing.T) {
	assert.NoErr(
		t,
		checkErr(errGettingServiceCatalog{Original: errors.New("testorig")}, "testorig"),
	)
}

func TestErrCreatingThirdPartyResource(t *testing.T) {
	assert.NoErr(
		t,
		checkErr(errCreatingThirdPartyResource{Original: errors.New("testorig")}, "testorig"),
	)
}
