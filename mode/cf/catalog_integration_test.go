// +build integration

package cf

import (
	"context"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/test-utils/k8s"
)

func TestCFCataloger(t *testing.T) {
	rootCtx := context.Background()
	httpCl := http.DefaultClient
	ctx, cancelFn := context.WithCancel(rootCtx)
	defer cancelFn()

	cataloger, _, err := GetComponents(ctx, httpCl)
	assert.NoErr(t, err)
	services, err := cataloger.List()
	assert.NoErr(t, err)
	// Compare to known results from cf-sample-broker...
	expectedServiceCount := 3
	expectedPlanCount := 4
	assert.Equal(t, len(services), expectedServiceCount, "service count")
	for _, service := range services {
		assert.Equal(t, len(service.Plans), expectedPlanCount, "plan count")
	}
	if err := k8s.DeleteNamespace(namespace); err != nil {
		t.Fatalf("Unexpected error deleteing namespace %s: %s", namespace, err)
	}
}
