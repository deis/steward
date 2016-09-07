package claim

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

type cmInterface struct {
	cm kcl.ConfigMapsInterface
}

func (c cmInterface) Get(name string) (*ServicePlanClaimWrapper, error) {
	cm, err := c.cm.Get(name)
	if err != nil {
		return nil, err
	}
	return servicePlanClaimWrapperFromConfigMap(cm)
}

func (c cmInterface) List(opts api.ListOptions) (*ServicePlanClaimsListWrapper, error) {
	cms, err := c.cm.List(opts)
	if err != nil {
		return nil, err
	}
	claims := make([]*ServicePlanClaimWrapper, len(cms.Items))
	for i, cm := range cms.Items {
		wr, err := servicePlanClaimWrapperFromConfigMap(&cm)
		if err != nil {
			return nil, err
		}
		claims[i] = wr
	}
	return &ServicePlanClaimsListWrapper{
		ResourceVersion: cms.ResourceVersion,
		Claims:          claims,
	}, nil
}

func (c cmInterface) Update(spc *ServicePlanClaimWrapper) (*ServicePlanClaimWrapper, error) {
	cm := &api.ConfigMap{
		Data:       spc.Claim.ToMap(),
		ObjectMeta: spc.ObjectMeta,
	}
	logger.Debugf("updating ConfigMap %s with data %s", cm.Name, cm.Data)
	newCM, err := c.cm.Update(cm)
	if err != nil {
		return nil, err
	}
	return servicePlanClaimWrapperFromConfigMap(newCM)
}

func (c cmInterface) Watch(opts api.ListOptions) Watcher {
	return newConfigMapWatcher(func() (watch.Interface, error) {
		return c.cm.Watch(opts)
	})
}
