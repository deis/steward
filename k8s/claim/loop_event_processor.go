package claim

import (
	"errors"
	"fmt"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
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

func newErrClaim(claim mode.ServicePlanClaim, err error) mode.ServicePlanClaim {
	claim.Status = mode.StatusFailed
	claim.StatusDescription = fmt.Sprintf("error: %s", err)
	return claim
}

func processEvent(
	ctx context.Context,
	evt *Event,
	cmNamespacer kcl.ConfigMapsNamespacer,
	catalog k8s.ServiceCatalogLookup,
	lifecycler mode.Lifecycler,
	claimCh chan<- mode.ServicePlanClaim,
) {

	claim := *evt.claim.Claim
	svc := catalog.Get(claim.ServiceID, claim.PlanID)
	if svc == nil {
		logger.Debugf("service %s, plan %s not found", claim.ServiceID, claim.PlanID)
		err := errNoSuchServiceAndPlan{
			svcID:  claim.ServiceID,
			planID: claim.PlanID,
		}
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}
	if claim.Action == mode.ActionCreate {
		logger.Debugf("creating service %s, plan %s from claim %s", svc.Info.Name, svc.Plan.Name, claim.ClaimID)
		processCreate(ctx, claim, svc, lifecycler, cmNamespacer, claimCh)
	} else if claim.Action == mode.ActionDelete {
		logger.Debugf("deleting service %s, plan %s from claim %s", svc.Info.Name, svc.Plan.Name, claim.ClaimID)
		processDelete(ctx, claim, svc, lifecycler, cmNamespacer, claimCh)
	}
}

func processCreate(
	ctx context.Context,
	claim mode.ServicePlanClaim,
	svc *k8s.ServiceCatalogEntry,
	lifecycler mode.Lifecycler,
	cmNamespacer kcl.ConfigMapsNamespacer,
	claimCh chan<- mode.ServicePlanClaim,
) {
	claim.Status = mode.StatusCreating
	select {
	case claimCh <- claim:
	case <-ctx.Done():
		return
	}

	// provision
	orgGUID := uuid.New()
	spaceGUID := uuid.New()
	instanceID := uuid.New()
	claim.Status = mode.StatusProvisioning
	claim.InstanceID = instanceID
	select {
	case claimCh <- claim:
	case <-ctx.Done():
		return
	}
	if _, err := lifecycler.Provision(instanceID, &mode.ProvisionRequest{
		OrganizationGUID: orgGUID,
		PlanID:           svc.Plan.ID,
		ServiceID:        svc.Info.ID,
		SpaceGUID:        spaceGUID,
		Parameters:       mode.JSONObject(map[string]string{}),
	}); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}

	// bind
	bindID := uuid.New()
	claim.Status = mode.StatusBinding
	claim.BindID = bindID
	select {
	case claimCh <- claim:
	case <-ctx.Done():
		return
	}
	bindRes, err := lifecycler.Bind(instanceID, bindID, &mode.BindRequest{})
	if err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}

	if _, err := cmNamespacer.ConfigMaps(claim.TargetNamespace).Create(&api.ConfigMap{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      claim.TargetName,
			Namespace: claim.TargetNamespace,
		},
		Data: bindRes.Creds,
	}); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}
	claim.Status = mode.StatusCreated
	select {
	case claimCh <- claim:
	case <-ctx.Done():
		return
	}
}

func processDelete(
	ctx context.Context,
	claim mode.ServicePlanClaim,
	svc *k8s.ServiceCatalogEntry,
	lifecycler mode.Lifecycler,
	cmNamespacer kcl.ConfigMapsNamespacer,
	claimCh chan<- mode.ServicePlanClaim,
) {
	claim.Status = mode.StatusDeleting
	select {
	case claimCh <- claim:
	case <-ctx.Done():
		return
	}

	// unbind
	instanceID := claim.InstanceID
	bindID := claim.BindID
	if instanceID == "" {
		select {
		case claimCh <- newErrClaim(claim, errMissingInstanceID):
		case <-ctx.Done():
			return
		}
	}
	if bindID == "" {
		select {
		case claimCh <- newErrClaim(claim, errMissingBindID):
		case <-ctx.Done():
			return
		}
	}

	claim.Status = mode.StatusUnbinding
	select {
	case claimCh <- claim:
	case <-ctx.Done():
		return
	}
	if err := lifecycler.Unbind(claim.ServiceID, claim.PlanID, instanceID, bindID); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}

	// deprovision
	if _, err := lifecycler.Deprovision(instanceID, claim.ServiceID, claim.PlanID); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}

	// delete configmap
	if err := cmNamespacer.ConfigMaps(claim.TargetNamespace).Delete(claim.TargetName); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-ctx.Done():
		}
		return
	}

	claim.Status = mode.StatusDeleted
	select {
	case claimCh <- claim:
	case <-ctx.Done():
	}
}
