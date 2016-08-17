package claim

import (
	"fmt"

	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/api"
)

// ServicePlanClaimWrapper is a wrapper for a ServicePlanClaim that also contains kubernetes-specific information
type ServicePlanClaimWrapper struct {
	Claim           *mode.ServicePlanClaim
	ResourceVersion string
	OriginalName    string
	Labels          map[string]string
}

func servicePlanClaimWrapperFromConfigMap(cm *api.ConfigMap) (*ServicePlanClaimWrapper, error) {
	claim, err := mode.ServicePlanClaimFromMap(cm.Data)
	if err != nil {
		return nil, err
	}
	return &ServicePlanClaimWrapper{
		Claim:           claim,
		ResourceVersion: cm.ResourceVersion,
		OriginalName:    cm.Name,
		Labels:          cm.Labels,
	}, nil
}

// String is the fmt.Stringer implementation
func (spc ServicePlanClaimWrapper) String() string {
	return fmt.Sprintf("%s (resource %s)", *spc.Claim, spc.ResourceVersion)
}

func (spc ServicePlanClaimWrapper) toConfigMap() *api.ConfigMap {
	return &api.ConfigMap{
		ObjectMeta: api.ObjectMeta{
			Name:            spc.OriginalName,
			Labels:          spc.Labels,
			ResourceVersion: spc.ResourceVersion,
		},
		Data: spc.Claim.ToMap(),
	}
}

// ServicePlanClaimsListWrapper is a wrapper for a list of ServicePlanClaims that also contains kubernetes-specific information.
type ServicePlanClaimsListWrapper struct {
	Claims          []*ServicePlanClaimWrapper
	ResourceVersion string
}
