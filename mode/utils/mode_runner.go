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
