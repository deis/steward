package utils

import (
	"context"
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
	"github.com/deis/steward/mode/cmd"
	"github.com/deis/steward/mode/helm"
	"google.golang.org/grpc"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/errors"
	"k8s.io/client-go/1.4/rest"
)

const (
	cfMode      = "cf"
	helmMode    = "helm"
	cmdMode     = "cmd"     // Official mode name
	commandMode = "command" // A config mistake I can see easily being made; should be forgiven
)

// Run publishes the underlying broker's service offerings to the catalog, then starts Steward's
// control loop in the specified mode. returns a function that should be used by the caller to clean up, should the run stop. this function will be
func Run(
	ctx context.Context,
	httpCl *http.Client,
	modeStr string,
	brokerName string,
	errCh chan<- error,
	namespaces []string,
) (func(), error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errGettingK8sClient{Original: err}
	}

	cleanupFunc := func() {}

	var cataloger mode.Cataloger
	var lifecycler *mode.Lifecycler
	// Get the right implementations of mode.Cataloger and mode.Lifecycler
	switch modeStr {
	case cfMode:
		cataloger, lifecycler, err = cf.GetComponents(ctx, httpCl)
		if err != nil {
			return nil, err
		}
	case helmMode:
		// ignore the returned connection. we're going to hold onto it for the whole lifespan of this steward, so we don't need to close it
		var conn *grpc.ClientConn
		conn, cataloger, lifecycler, err = helm.GetComponents(ctx, httpCl, k8sClient)
		if err != nil {
			return nil, err
		}
		cleanupFunc = func() { conn.Close() }
	case cmdMode:
		fallthrough
	case commandMode:
		cataloger, lifecycler, err = cmd.GetComponents(k8sClient)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errUnrecognizedMode{mode: modeStr}
	}

	// Everything else does not vary by mode...

	// Create Service Catalog 3PR without bombing out if it already exists.
	extensions := k8sClient.Extensions()
	tpr := extensions.ThirdPartyResources()

	_, err = tpr.Create(k8s.ServiceCatalog3PR)
	if err != nil && !errors.IsAlreadyExists(err) {
		return nil, errCreatingThirdPartyResource{Original: err}
	}

	catalogInteractor := k8s.NewK8sServiceCatalogInteractor(k8sClient.CoreClient.RESTClient)
	published, err := publishCatalog(brokerName, cataloger, catalogInteractor)
	if err != nil {
		return nil, errPublishingServiceCatalog{Original: err}
	}
	logger.Infof("published %d entries into the service catalog", len(published))

	evtNamespacer := claim.NewConfigMapsInteractorNamespacer(k8sClient)
	lookup, err := k8s.FetchServiceCatalogLookup(catalogInteractor)
	if err != nil {
		return nil, errGettingServiceCatalogLookupTable{Original: err}
	}
	logger.Infof("created service catalog lookup with %d items", lookup.Len())
	claim.StartControlLoops(ctx, evtNamespacer, k8sClient, *lookup, lifecycler, namespaces, errCh)

	return cleanupFunc, nil
}

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
