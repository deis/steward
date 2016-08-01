package k8s

import (
	"encoding/json"
	// "errors"
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	kcl "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	servicePlanClaims3PRName = "serviceplanclaims"
)

// ServicePlanClaim represents the JSON data an application should put into a Kubernetes third party resource to tell steward that it wants to provision a service. See https://github.com/deis/steward/blob/master/DATA_STRUCTURES.md#serviceplanclaim for more detail
//
// This type implements (k8s.io/kubernetes/pkg/runtime).Object
type ServicePlanClaim struct {
	ServiceProvider   string `json:"service_provider"`
	ServicePlan       string `json:"service_plan"`
	ClaimID           string `json:"claim_id"`
	Action            string `json:"action"`
	Status            string `json:"status"`
	StatusDescription string `json:"status_description"`
}

// GetObjectKind is the (k8s.io/kubernetes/pkg/runtime).Object interface implementation
func (s ServicePlanClaim) GetObjectKind() kcl.ObjectKind {
	return &kcl.TypeMeta{Kind: "serviceplanclaim", APIVersion: resourceAPIVersion(apiVersionV1)}
}

func (s ServicePlanClaim) String() string {
	return fmt.Sprintf(
		"%s (%s)\nclaimID: %s\naction: %s\nstatus: %s\nstatusDescription: %s\n",
		s.ServiceProvider,
		s.ServicePlan,
		s.ClaimID,
		s.Action,
		s.Status,
		s.StatusDescription,
	)
}

type servicePlanClaim3PRWrapper struct {
	api.ObjectMeta       `json:"metadata,omitempty"`
	unversioned.TypeMeta `json:",inline"`
	*ServicePlanClaim    `json:",inline"`
}

// ServicePlanClaims3PRWrapper is the data structure returned by the GET call for all of the service plan claims
type servicePlanClaims3PRWrapper struct {
	api.ObjectMeta       `json:"metadata,omitempty"`
	unversioned.TypeMeta `json:",inline"`
	Items                []*ServicePlanClaim `json:"items"`
}

func getServicePlanClaimsAbsPath(namespace string) []string {
	return []string{"apis", "steward.deis.com", "v1", "namespaces", namespace, servicePlanClaims3PRName}
}

// get service plan claims in namespaces from k8s using cl
func getServicePlanClaims(cl *restclient.RESTClient, namespace string) ([]*ServicePlanClaim, error) {
	req := cl.Get().AbsPath(getServicePlanClaimsAbsPath(namespace)...)
	resBody, err := req.DoRaw()
	if err != nil {
		return nil, err
	}

	lst := new(servicePlanClaims3PRWrapper)
	if err := json.Unmarshal(resBody, lst); err != nil {
		return nil, err
	}
	return lst.Items, nil
}

// launches a new goroutine to query the service plan claims endpoint in namespace using cl. this func maintains a cache and, after each claim, returns those that haven't been returned already. pauses sleepDur between queries for claims, and on each iteration either sends claims on the first returned channel or errors on the second. stops the internal goroutine if the first channel is closed
func watchServicePlanClaims(
	cl *restclient.RESTClient,
	namespace string,
	sleepDur time.Duration,
) (chan<- struct{}, <-chan *ServicePlanClaim, <-chan error) {
	spcCh := make(chan *ServicePlanClaim)
	errCh := make(chan error)
	stopCh := make(chan struct{})
	go func() {
		cache := make(map[string]struct{})
		for {
			select {
			case <-stopCh:
				return
			default:
			}
			claimsList, err := getServicePlanClaims(cl, namespace)
			if err != nil {
				errCh <- err
				time.Sleep(sleepDur)
				continue
			}
			for _, claim := range claimsList {
				if _, ok := cache[claim.String()]; ok {
					continue
				}
				spcCh <- claim
				cache[claim.String()] = struct{}{}
			}
			time.Sleep(sleepDur)
		}
	}()
	return stopCh, spcCh, errCh
}
