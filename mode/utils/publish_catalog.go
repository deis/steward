package utils

import (
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/pkg/api"
)

// Does the following:
//
//	1. fetches the service catalog from the backing broker
//	2. checks the 3pr for already-existing entries, and errors if one already exists
//	3. if none error-ed in #2, publishes 3prs for all of the catalog entries
//
// returns all of the entries it wrote into the catalog, or an error
func publishCatalog(
	brokerName string,
	cataloger mode.Cataloger,
	catalogEntries k8s.ServiceCatalogInteractor,
) ([]*k8s.ServiceCatalogEntry, error) {

	services, err := cataloger.List()
	if err != nil {
		return nil, errGettingServiceCatalog{Original: err}
	}

	published := []*k8s.ServiceCatalogEntry{}
	// Write all entries from cf catalog to 3prs
	for _, service := range services {
		for _, plan := range service.Plans {
			entry := k8s.NewServiceCatalogEntry(brokerName, api.ObjectMeta{}, service.ServiceInfo, plan)
			if _, err := catalogEntries.Create(entry); err != nil {
				logger.Errorf(
					"error publishing catalog entry (svc_name, plan_name) = (%s, %s) (%s), continuing",
					entry.Info.Name,
					entry.Plan.Name,
					err,
				)
				continue
			}
			published = append(published, entry)
		}
	}

	return published, nil
}
