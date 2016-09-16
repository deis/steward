package claim

import (
	"context"

	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	"k8s.io/client-go/1.4/pkg/api"
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/watch"
)

type cmInterface struct {
	cm v1.ConfigMapInterface
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
	cm := &v1types.ConfigMap{
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

func (c cmInterface) Watch(ctx context.Context, opts api.ListOptions) Watcher {
	return newConfigMapWatcher(ctx, func() (watch.Interface, error) {
		return c.cm.Watch(opts)
	})
}
