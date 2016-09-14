package helm

import (
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// look up the ConfigMap using cmNamespacer for each item in cmInfos. return any errors getting config map or calling fn
func rangeConfigMaps(
	cmNamespacer kcl.ConfigMapsNamespacer,
	cmInfos []cmNamespaceAndName,
	fn func(*api.ConfigMap) error) error {
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
