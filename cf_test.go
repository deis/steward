package main

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/fake"
)

func TestPublishCatalog(t *testing.T) {
	cataloger := fake.Cataloger{
		Services: []*mode.Service{
			&mode.Service{
				ServiceInfo: mode.ServiceInfo{
					Name:        "service 1",
					ID:          "service1",
					Description: "this is service 1",
				},
				Plans: []mode.ServicePlan{
					mode.ServicePlan{
						ID:          "plan1",
						Name:        "plan 1",
						Description: "this is plan 1",
					},
					mode.ServicePlan{
						ID:          "plan2",
						Name:        "plan 2",
						Description: "this is plan 2",
					},
				},
			},
		},
	}
	catalogEntries := k8s.FakeServiceCatalogInteractor{
		ListRet: &k8s.ServiceCatalogEntryList{Items: nil},
	}
	entries, err := publishCatalog(cataloger, &catalogEntries)
	assert.NoErr(t, err)
	expectedNumEntries := 0
	for _, svc := range cataloger.Services {
		expectedNumEntries += len(svc.Plans)
	}
	assert.Equal(t, len(entries), expectedNumEntries, "number of entries published")
	assert.Equal(t, len(catalogEntries.Created), expectedNumEntries, "number of created entries")

	catalogEntries = k8s.FakeServiceCatalogInteractor{
		ListRet:   &k8s.ServiceCatalogEntryList{Items: nil},
		CreateErr: errors.New("test error"),
	}
	entries, err = publishCatalog(cataloger, &catalogEntries)
	assert.NoErr(t, err)
	assert.Equal(t, len(entries), 0, "number of entries published")
	assert.Equal(t, len(catalogEntries.Created), expectedNumEntries, "number of entries created")
}
