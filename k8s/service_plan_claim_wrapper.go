package k8s

import (
	"fmt"

	"k8s.io/client-go/1.4/pkg/api/v1"
)

// ServicePlanClaimWrapper is a wrapper for a ServicePlanClaim that also contains kubernetes-specific information
type ServicePlanClaimWrapper struct {
	ObjectMeta v1.ObjectMeta
	Claim      *ServicePlanClaim
}

// ServicePlanClaimWrapperFromConfigMap parses a ServicePlanClaim from cm and returns the wrapper representation of it. Returns nil and an error if the config map was malformed
func ServicePlanClaimWrapperFromConfigMap(cm *v1.ConfigMap) (*ServicePlanClaimWrapper, error) {
	claim, err := ServicePlanClaimFromMap(cm.Data)
	if err != nil {
		return nil, err
	}
	return &ServicePlanClaimWrapper{
		Claim: claim,
		ObjectMeta: v1.ObjectMeta{
			ResourceVersion: cm.ResourceVersion,
			Name:            cm.Name,
			Namespace:       cm.Namespace,
			Labels:          cm.Labels,
		},
	}, nil
}

// String is the fmt.Stringer implementation
func (spc ServicePlanClaimWrapper) String() string {
	return fmt.Sprintf("%s (resource %s)", *spc.Claim, spc.ObjectMeta.ResourceVersion)
}

// ToConfigMap converts spc to a ConfigMap that ServicePlanClaimFromMap will decode back into spc
func (spc ServicePlanClaimWrapper) ToConfigMap() *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: spc.ObjectMeta,
		Data:       spc.Claim.ToMap(),
	}
}

// ServicePlanClaimsListWrapper is a wrapper for a list of ServicePlanClaims that also contains kubernetes-specific information.
type ServicePlanClaimsListWrapper struct {
	Claims          []*ServicePlanClaimWrapper
	ResourceVersion string
}
