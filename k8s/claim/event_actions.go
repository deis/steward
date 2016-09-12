package claim

import (
	"context"
	"errors"
	"fmt"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	errMissingInstanceID = errors.New("missing instance ID")
	errMissingBindID     = errors.New("missing bind ID")
)

type errNoSuchServiceAndPlan struct {
	svcID  string
	planID string
}

func (e errNoSuchServiceAndPlan) Error() string {
	return fmt.Sprintf("no such service and plan. service ID = %s, plan ID = %s", e.svcID, e.planID)
}

func isNoSuchServiceAndPlanErr(e error) bool {
	_, ok := e.(errNoSuchServiceAndPlan)
	return ok
}

func getService(claim mode.ServicePlanClaim, catalog k8s.ServiceCatalogLookup) (*k8s.ServiceCatalogEntry, error) {
	svc := catalog.Get(claim.ServiceID, claim.PlanID)
	if svc == nil {
		logger.Debugf("service %s, plan %s not found", claim.ServiceID, claim.PlanID)
		return nil, errNoSuchServiceAndPlan{
			svcID:  claim.ServiceID,
			planID: claim.PlanID,
		}
	}
	return svc, nil
}

func processProvision(
	ctx context.Context,
	evt *Event,
	secretsNamespacer kcl.SecretsNamespacer,
	catalogLookup k8s.ServiceCatalogLookup,
	lifecycler *mode.Lifecycler,
	claimCh chan<- state.Update,
) {

	logger.Debugf("processProvision for event %s", evt.claim.Claim.ToMap())

	claim := *evt.claim.Claim

	svc, err := getService(claim, catalogLookup)
	if err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	select {
	case claimCh <- state.StatusUpdate(mode.StatusProvisioning):
	case <-ctx.Done():
		return
	}
	orgGUID := uuid.New()
	spaceGUID := uuid.New()
	instanceID := uuid.New()
	provisionResp, err := lifecycler.Provision(instanceID, &mode.ProvisionRequest{
		OrganizationGUID: orgGUID,
		PlanID:           svc.Plan.ID,
		ServiceID:        svc.Info.ID,
		SpaceGUID:        spaceGUID,
		Parameters:       mode.JSONObject(map[string]string{}),
	})
	if err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}
	select {
	case claimCh <- state.FullUpdate(mode.StatusProvisioned, "", instanceID, "", provisionResp.Extra):
	case <-ctx.Done():
		return
	}
}

func processBind(
	ctx context.Context,
	evt *Event,
	secretsNamespacer kcl.SecretsNamespacer,
	catalogLookup k8s.ServiceCatalogLookup,
	lifecycler *mode.Lifecycler,
	claimCh chan<- state.Update,
) {

	logger.Debugf("processBind for event %s", evt.claim.Claim.ToMap())

	claim := *evt.claim.Claim
	claimWrapper := *evt.claim
	if _, err := getService(claim, catalogLookup); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	select {
	case claimCh <- state.StatusUpdate(mode.StatusBinding):
	case <-ctx.Done():
		return
	}

	instanceID := claim.InstanceID
	if instanceID == "" {
		select {
		case claimCh <- state.ErrUpdate(errMissingInstanceID):
		case <-ctx.Done():
		}
		return
	}

	bindID := uuid.New()
	bindRes, err := lifecycler.Bind(instanceID, bindID, &mode.BindRequest{
		ServiceID:  claim.ServiceID,
		PlanID:     claim.PlanID,
		Parameters: mode.JSONObject(map[string]string{}),
	})
	if err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	credBytes := make(map[string][]byte, len(bindRes.Creds))
	for k, v := range bindRes.Creds {
		credBytes[k] = []byte(v)
	}

	if _, err := secretsNamespacer.Secrets(claimWrapper.ObjectMeta.Namespace).Create(&api.Secret{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      claim.TargetName,
			Namespace: claimWrapper.ObjectMeta.Namespace,
		},
		Data: credBytes,
	}); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}
	select {
	case claimCh <- state.FullUpdate(mode.StatusBound, "", instanceID, bindID, mode.EmptyJSONObject()):
	case <-ctx.Done():
		return
	}
}

func processUnbind(
	ctx context.Context,
	evt *Event,
	secretsNamespacer kcl.SecretsNamespacer,
	catalogLookup k8s.ServiceCatalogLookup,
	lifecycler *mode.Lifecycler,
	claimCh chan<- state.Update,
) {

	logger.Debugf("processUnbind for event %s", evt.claim.Claim.ToMap())

	claimWrapper := evt.claim
	claim := *evt.claim.Claim
	if _, err := getService(claim, catalogLookup); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	select {
	case claimCh <- state.StatusUpdate(mode.StatusUnbinding):
	case <-ctx.Done():
		return
	}

	instanceID := claim.InstanceID
	bindID := claim.BindID
	if instanceID == "" {
		select {
		case claimCh <- state.ErrUpdate(errMissingInstanceID):
		case <-ctx.Done():
			return
		}
	}
	if bindID == "" {
		select {
		case claimCh <- state.ErrUpdate(errMissingBindID):
		case <-ctx.Done():
			return
		}
	}

	if err := lifecycler.Unbind(instanceID, bindID, &mode.UnbindRequest{
		ServiceID: claim.ServiceID,
		PlanID:    claim.PlanID,
	}); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	// delete secret
	if err := secretsNamespacer.Secrets(claimWrapper.ObjectMeta.Namespace).Delete(claim.TargetName); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	select {
	case claimCh <- state.StatusUpdate(mode.StatusUnbound):
	case <-ctx.Done():
		return
	}
}

func processDeprovision(
	ctx context.Context,
	evt *Event,
	secretsNamespacer kcl.SecretsNamespacer,
	catalogLookup k8s.ServiceCatalogLookup,
	lifecycler *mode.Lifecycler,
	claimCh chan<- state.Update,
) {

	logger.Debugf("processDeprovision for event %s", evt.claim.Claim.ToMap())

	claim := *evt.claim.Claim
	if _, err := getService(claim, catalogLookup); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}

	select {
	case claimCh <- state.StatusUpdate(mode.StatusDeprovisioning):
	case <-ctx.Done():
		return
	}

	instanceID := claim.InstanceID
	if instanceID == "" {
		select {
		case claimCh <- state.ErrUpdate(errMissingInstanceID):
		case <-ctx.Done():
		}
		return
	}

	// deprovision
	deprovisionReq := &mode.DeprovisionRequest{
		ServiceID:  claim.ServiceID,
		PlanID:     claim.PlanID,
		Parameters: evt.claim.Claim.Extra,
	}
	if _, err := lifecycler.Deprovision(instanceID, deprovisionReq); err != nil {
		select {
		case claimCh <- state.ErrUpdate(err):
		case <-ctx.Done():
		}
		return
	}
	claim.Status = mode.StatusDeprovisioned.String()
	select {
	case claimCh <- state.StatusUpdate(mode.StatusDeprovisioned):
	case <-ctx.Done():
		return
	}
}
