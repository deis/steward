package helm

import (
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
)

// look up the ConfigMap using cmNamespacer for each item in cmInfos. return any errors getting config map or calling fn
func rangeConfigMaps(
	cmNamespacer v1.ConfigMapsGetter,
	cmInfos []cmNamespaceAndName,
	fn func(*v1types.ConfigMap) error) error {
	for _, cmInfo := range cmInfos {
		logger.Debugf("getting config map %s/%s", cmInfo.Namespace, cmInfo.Name)
		cm, err := cmNamespacer.ConfigMaps(cmInfo.Namespace).Get(cmInfo.Name)
		if err != nil {
			logger.Debugf("no such ConfigMap %s/%s", cmInfo.Namespace, cmInfo.Name)
			return err
		}
		if err := fn(cm); err != nil {
			return err
		}
	}
	return nil
}
