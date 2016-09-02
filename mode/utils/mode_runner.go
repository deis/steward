package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim"
	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	cfMode   = "cf"
	helmMode = "helm"
)

// Run publishes the underlying broker's service offerings to the catalog, then starts Steward's
// control loop in the specified mode.
func Run(
	ctx context.Context,
	httpCl *http.Client,
	modeStr string,
	errCh chan<- error, namespaces []string) error {

	k8sClient, err := kcl.NewInCluster()
	if err != nil {
		return errGettingK8sClient{Original: err}
	}

	var cataloger mode.Cataloger
	var lifecycler *mode.Lifecycler
	// Get the right implementations of mode.Cataloger and mode.Lifecycler
	switch modeStr {
	case cfMode:
		var err error
		cataloger, lifecycler, err = getCfModeComponents(ctx)
		if err != nil {
			return err
		}
	case helmMode:
		var err error
		cataloger, lifecycler, err = getHelmModeComponents(ctx, httpCl, k8sClient)
		if err != nil {
			return err
		}
	default:
		return errUnrecognizedMode{mode: modeStr}
	}

	// Everything else does not vary by mode...

	catalogInteractor := k8s.NewK8sServiceCatalogInteractor(k8sClient.RESTClient)
	published, err := publishCatalog(cataloger, catalogInteractor)
	if err != nil {
		return errPublishingServiceCatalog{Original: err}
	}
	logger.Infof("published %d entries into the service catalog", len(published))

	evtNamespacer := claim.NewConfigMapsInteractorNamespacer(k8sClient)
	lookup, err := k8s.FetchServiceCatalogLookup(catalogInteractor)
	if err != nil {
		return errGettingServiceCatalogLookupTable{Original: err}
	}
	logger.Infof("created service catalog lookup with %d items", lookup.Len())
	claim.StartControlLoops(ctx, evtNamespacer, k8sClient, *lookup, lifecycler, namespaces, errCh)

	return nil
}

// Does the following:
//
//	1. fetches the service catalog from the backing broker
//	2. checks the 3pr for already-existing entries, and errors if one already exists
//	3. if none error-ed in #2, publishes 3prs for all of the catalog entries
//
// returns all of the entries it wrote into the catalog, or an error
func publishCatalog(
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
			descr := fmt.Sprintf("%s (%s)", service.Description, plan.Description)
			entry := k8s.NewServiceCatalogEntry(descr, api.ObjectMeta{}, service.ServiceInfo, plan)
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
