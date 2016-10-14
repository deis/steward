package k8s

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"k8s.io/client-go/1.4/pkg/api"
)

func TestFetchServiceCatalogLookup(t *testing.T) {
	iface := &FakeServiceCatalogInteractor{
		ListRet: &ServiceCatalogEntryList{},
	}
	lookup, err := FetchServiceCatalogLookup(iface)
	assert.NoErr(t, err)
	assert.Equal(t, lookup.Len(), 0, "number of items in the catalog")
	iface.ListRet.Items = []*ServiceCatalogEntry{
		NewServiceCatalogEntry(
			"testBroker1",
			api.ObjectMeta{},
			mode.ServiceInfo{
				ID:   uuid.New(),
				Name: "testSvc1",
			},
			mode.ServicePlan{
				ID:   uuid.New(),
				Name: "testPlan1",
			},
		),
		NewServiceCatalogEntry(
			"testBroker2",
			api.ObjectMeta{},
			mode.ServiceInfo{
				ID:   uuid.New(),
				Name: "testSvc2",
			},
			mode.ServicePlan{
				ID:   uuid.New(),
				Name: "testPlan2",
			},
		),
	}
	lookup, err = FetchServiceCatalogLookup(iface)
	assert.NoErr(t, err)
	assert.Equal(t, lookup.Len(), len(iface.ListRet.Items), "number of items in the catalog")
	for _, entry := range iface.ListRet.Items {
		assert.Equal(
			t,
			lookup.Get(entry.Info.ID, entry.Plan.ID),
			entry,
			fmt.Sprintf("service %s, plan %s", entry.Info.ID, entry.Plan.ID),
		)
	}
}

func TestServiceCatalogLookupCatalogKey(t *testing.T) {
	svcID := uuid.New()
	planID := uuid.New()

	key := catalogKey(svcID, planID)
	assert.Equal(t, key, fmt.Sprintf("%s-%s", svcID, planID), "catalog key")
}

func TestServiceCatalogLookupGetSet(t *testing.T) {
	initial := []*ServiceCatalogEntry{
		&ServiceCatalogEntry{
			Info: mode.ServiceInfo{ID: "testsvc1"},
			Plan: mode.ServicePlan{ID: "testplan1"},
		},
		&ServiceCatalogEntry{
			Info: mode.ServiceInfo{ID: "testsvc2"},
			Plan: mode.ServicePlan{ID: "testplan2"},
		},
	}
	lookup := NewServiceCatalogLookup(initial)
	for _, entry := range initial {
		fetched := lookup.Get(entry.Info.ID, entry.Plan.ID)
		if fetched == nil {
			t.Fatalf("service expected but not found for service ID %s, plan ID %s", entry.Info.ID, entry.Plan.ID)
		}
		assert.Equal(t, fetched.Info.ID, entry.Info.ID, "service ID")
		assert.Equal(t, fetched.Plan.ID, entry.Plan.ID, "plan ID")
	}
	newEntry := &ServiceCatalogEntry{
		Info: mode.ServiceInfo{ID: "testsvc3"},
		Plan: mode.ServicePlan{ID: "testplan3"},
	}
	lookup.Set(newEntry)
	fetched := lookup.Get(newEntry.Info.ID, newEntry.Plan.ID)
	if fetched == nil {
		t.Fatalf("entry with service ID %s, plan ID %s not found after it was set", newEntry.Info.ID, newEntry.Plan.ID)
	}
	assert.Equal(t, fetched.Info.ID, newEntry.Info.ID, "service ID")
	assert.Equal(t, fetched.Plan.ID, newEntry.Plan.ID, "plan ID")
}
