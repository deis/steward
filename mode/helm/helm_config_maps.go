package helm

import (
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
)

// look up the ConfigMap using cmNamespacer for each item in cmInfos. return any errors getting config map or calling fn
func rangeConfigMaps(
	cmIface v1.ConfigMapInterface,
	cmNames []string,
	fn func(*v1types.ConfigMap) error) error {
	for _, cmName := range cmNames {
		logger.Debugf("getting config map %s", cmName)
		cm, err := cmIface.Get(cmName)
		if err != nil {
			logger.Debugf("no such ConfigMap %s", cmName)
			return err
		}
		if err := fn(cm); err != nil {
			return err
		}
	}
	return nil
}
