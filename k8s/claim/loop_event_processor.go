package claim

import (
	"errors"
	"fmt"

	"github.com/deis/steward/k8s"
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

func newErrClaim(claim mode.ServicePlanClaim, err error) mode.ServicePlanClaim {
	claim.Status = mode.StatusFailed
	claim.StatusDescription = fmt.Sprintf("error: %s", err)
	return claim
}

func processEvent(
	evt *Event,
	cmNamespacer kcl.ConfigMapsNamespacer,
	catalog k8s.ServiceCatalogLookup,
	lifecycler mode.Lifecycler,
	claimCh chan<- mode.ServicePlanClaim,
	doneCh <-chan struct{},
) {

	defer close(claimCh)

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
		case <-doneCh:
		}
		return
	}
	if claim.Action == mode.ActionCreate {
		logger.Debugf("creating service %s, plan %s from claim %s", svc.Info.Name, svc.Plan.Name, claim.ClaimID)
		processCreate(claim, svc, lifecycler, cmNamespacer, claimCh, doneCh)
	} else if claim.Action == mode.ActionDelete {
		logger.Debugf("deleting service %s, plan %s from claim %s", svc.Info.Name, svc.Plan.Name, claim.ClaimID)
		processDelete(claim, svc, lifecycler, cmNamespacer, claimCh, doneCh)
	}
}

func processCreate(
	claim mode.ServicePlanClaim,
	svc *k8s.ServiceCatalogEntry,
	lifecycler mode.Lifecycler,
	cmNamespacer kcl.ConfigMapsNamespacer,
	claimCh chan<- mode.ServicePlanClaim,
	doneCh <-chan struct{},
) {
	claim.Status = mode.StatusCreating
	select {
	case claimCh <- claim:
	case <-doneCh:
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
	case <-doneCh:
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
		case <-doneCh:
		}
		return
	}

	// bind
	bindID := uuid.New()
	claim.Status = mode.StatusBinding
	claim.BindID = bindID
	select {
	case claimCh <- claim:
	case <-doneCh:
		return
	}
	bindRes, err := lifecycler.Bind(instanceID, bindID, &mode.BindRequest{})
	if err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-doneCh:
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
		case <-doneCh:
		}
		return
	}
	claim.Status = mode.StatusCreated
	select {
	case claimCh <- claim:
	case <-doneCh:
		return
	}
}

func processDelete(
	claim mode.ServicePlanClaim,
	svc *k8s.ServiceCatalogEntry,
	lifecycler mode.Lifecycler,
	cmNamespacer kcl.ConfigMapsNamespacer,
	claimCh chan<- mode.ServicePlanClaim,
	doneCh <-chan struct{},
) {
	claim.Status = mode.StatusDeleting
	select {
	case claimCh <- claim:
	case <-doneCh:
		return
	}

	// unbind
	instanceID := claim.InstanceID
	bindID := claim.BindID
	if instanceID == "" {
		select {
		case claimCh <- newErrClaim(claim, errMissingInstanceID):
		case <-doneCh:
			return
		}
	}
	if bindID == "" {
		select {
		case claimCh <- newErrClaim(claim, errMissingBindID):
		case <-doneCh:
			return
		}
	}

	claim.Status = mode.StatusUnbinding
	select {
	case claimCh <- claim:
	case <-doneCh:
		return
	}
	if err := lifecycler.Unbind(claim.ServiceID, claim.PlanID, instanceID, bindID); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-doneCh:
		}
		return
	}

	// deprovision
	if _, err := lifecycler.Deprovision(instanceID, claim.ServiceID, claim.PlanID); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-doneCh:
		}
		return
	}

	// delete configmap
	if err := cmNamespacer.ConfigMaps(claim.TargetNamespace).Delete(claim.TargetName); err != nil {
		select {
		case claimCh <- newErrClaim(claim, err):
		case <-doneCh:
		}
		return
	}

	claim.Status = mode.StatusDeleted
	select {
	case claimCh <- claim:
	case <-doneCh:
	}
}
