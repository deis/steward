package main

import (
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
	"github.com/deis/steward/web"
	"github.com/deis/steward/web/brokerapi"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func runCFMode(
	apiServerHostStr string,
	frontendAuth *web.BasicAuth,
	cl *kcl.Client,
	errCh chan<- error,
	namespaces []string,
) error {

	cfCfg, err := cf.GetConfig()
	if err != nil {
		logger.Criticalf("error getting CloudFoundry broker config (%s)", err)
		return err
	}
	logger.Infof(
		"starting in CloudFoundry mode with hostname %s, port %d and username %s",
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
	)
	cfClient := cf.NewRESTClient(
		http.DefaultClient,
		cfCfg.Scheme,
		cfCfg.Hostname,
		cfCfg.Port,
		cfCfg.Username,
		cfCfg.Password,
	)

	catalogInteractor := k8s.NewK8sServiceCatalogInteractor(cl.RESTClient)
	cataloger := cf.NewCataloger(cfClient)
	lifecycler := cf.NewLifecycler(cfClient)

	published, err := publishCatalog(cataloger, catalogInteractor)
	if err != nil {
		logger.Criticalf("error publishing the cloud foundry service catalog (%s)", err)
		return err
	}

	logger.Infof("published %d entries into the catalog", len(published))
	go runBrokerAPI(cataloger, lifecycler, frontendAuth, apiServerHostStr, errCh, cl)

	evtNamespacer := claim.NewConfigMapsInteractorNamespacer(cl)
	lookup, err := k8s.FetchServiceCatalogLookup(catalogInteractor)
	if err != nil {
		logger.Criticalf("error fetching the service catalog lookup table (%s)", err)
		return err
	}
	logger.Infof("created service catalog lookup with %d items", lookup.Len())
	go claim.StartControlLoops(evtNamespacer, cl, *lookup, lifecycler, namespaces, errCh)

	return nil
}

// does the following:
//
//	1. fetches the service catalog from cloud foundry
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
		logger.Debugf("error getting CF catalog (%s)", err)
		return nil, err
	}

	published := []*k8s.ServiceCatalogEntry{}
	// write all entries from cf catalog to 3pr
	for _, service := range services {
		for _, plan := range service.Plans {
			entry := k8s.NewServiceCatalogEntry(api.ObjectMeta{}, service.ServiceInfo, plan)
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

func runBrokerAPI(
	cataloger mode.Cataloger,
	lifecycler mode.Lifecycler,
	frontendAuth *web.BasicAuth,
	hostStr string,
	errCh chan<- error,
	cmNamespacer kcl.ConfigMapsNamespacer,
) {

	logger.Infof("starting CF broker API server on %s", hostStr)
	hdl := brokerapi.Handler(cataloger, lifecycler, frontendAuth, cmNamespacer)
	if err := http.ListenAndServe(hostStr, hdl); err != nil {
		errCh <- err
	}
}
