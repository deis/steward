package main

import (
	"fmt"
	"net/http"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
	"github.com/deis/steward/web"
	"github.com/deis/steward/web/brokerapi"
	"github.com/juju/loggo"
	"k8s.io/kubernetes/pkg/client/restclient"
)

type errServiceAlreadyPublished struct {
	entry *k8s.ServiceCatalogEntry
}

func (e errServiceAlreadyPublished) Error() string {
	return fmt.Sprintf("duplicate service catalog entry: %s", e.entry.ResourceName())
}

func runCFMode(
	logger loggo.Logger,
	apiServerHostStr string,
	frontendAuth *web.BasicAuth,
	cl *restclient.RESTClient,
	errCh chan<- error,
	cmCreatorDeleter k8s.ConfigMapCreatorDeleter,
	secCreatorDeleter k8s.SecretCreatorDeleter,
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
	provisioner := cf.NewProvisioner(logger, cfClient)
	binder := cf.NewBinder(logger, cfClient)
	unbinder := cf.NewUnbinder(logger, cfClient)
	cataloger := cf.NewCataloger(logger, cfClient)

	published, err := publishCatalog(logger, cataloger, cl)
	if err != nil {
		logger.Criticalf("error publishing the cloud foundry service catalog (%s)", err)
		return err
	}
	logger.Infof("published %d entries into the catalog", len(published))
	for _, pub := range published {
		logger.Debugf("%s", pub.Info.Name)
	}
	go runBrokerAPI(
		logger,
		cataloger,
		provisioner,
		binder,
		unbinder,
		frontendAuth,
		apiServerHostStr,
		errCh,
		cmCreatorDeleter,
		secCreatorDeleter,
	)
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
	logger loggo.Logger,
	cataloger mode.Cataloger,
	restCl *restclient.RESTClient,
) ([]*k8s.ServiceCatalogEntry, error) {

	services, err := cataloger.List()
	if err != nil {
		logger.Debugf("error getting CF catalog (%s)", err)
		return nil, err
	}

	// get existing catalog from 3pr
	catalogEntries, err := k8s.GetServiceCatalogEntries(logger, restCl)
	if err != nil {
		logger.Debugf("error getting existing service catalog entries from k8s (%s)", err)
		return nil, err
	}
	// create in-mem lookup table from 3pr catalog, check for duplicate entries in cf catalog
	lookup := k8s.NewServiceCatalogLookup(catalogEntries)
	for _, service := range services {
		for _, plan := range service.Plans {
			entry := &k8s.ServiceCatalogEntry{Info: service.ServiceInfo, Plan: plan}
			if lookup.Get(entry) != nil {
				return nil, errServiceAlreadyPublished{entry: entry}
			}
			lookup.Set(entry)
		}
	}
	published := []*k8s.ServiceCatalogEntry{}
	// write all entries from cf catalog to 3pr
	for _, service := range services {
		for _, plan := range service.Plans {
			entry := &k8s.ServiceCatalogEntry{Info: service.ServiceInfo, Plan: plan}
			if err := k8s.PublishServiceCatalogEntry(restCl, entry); err != nil {
				logger.Errorf("error publishing catalog entry %s (%s), continuing", entry.ResourceName(), err)
				continue
			}
			published = append(published, entry)
		}
	}

	return published, nil
}

func runBrokerAPI(
	logger loggo.Logger,
	cataloger mode.Cataloger,
	provisioner mode.Provisioner,
	binder mode.Binder,
	unbinder mode.Unbinder,
	frontendAuth *web.BasicAuth,
	hostStr string,
	errCh chan<- error,
	cmCreatorDeleter k8s.ConfigMapCreatorDeleter,
	secCreatorDeleter k8s.SecretCreatorDeleter,
) {

	logger.Infof("starting CF broker API server on %s", hostStr)
	hdl := brokerapi.Handler(
		logger,
		cataloger,
		provisioner,
		binder,
		unbinder,
		frontendAuth,
		cmCreatorDeleter,
		secCreatorDeleter,
	)
	if err := http.ListenAndServe(hostStr, hdl); err != nil {
		errCh <- err
	}
}
