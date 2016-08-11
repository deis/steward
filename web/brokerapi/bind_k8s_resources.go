package brokerapi

import (
	"fmt"
	"time"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"github.com/juju/loggo"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	configMapKind = "ConfigMap"
)

func getObjectMeta(namespace, name string) api.ObjectMeta {
	standardLabels := map[string]string{
		"created-by": "steward",
		"created-at": fmt.Sprintf("%d", time.Now().Unix()),
	}
	return api.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    standardLabels,
	}
}

// writes everything in creds to a config map with the given name and namespace. returns error if the create failed
func writeToKubernetes(
	logger loggo.Logger,
	namespace,
	name string,
	creds mode.JSONObject,
	cmCreator k8s.ConfigMapCreator,
) error {

	configMap := &api.ConfigMap{
		TypeMeta:   unversioned.TypeMeta{Kind: configMapKind},
		ObjectMeta: getObjectMeta(namespace, name),
		Data:       creds.Base64EncodedVals(),
	}

	logger.Debugf("creating config map with bind credentials %+v", *configMap)
	if _, err := cmCreator.Create(namespace, configMap); err != nil {
		return err
	}
	return nil
}

func deleteFromKubernetes(namespace, name string, configMapDeleter k8s.ConfigMapDeleter) error {
	return configMapDeleter.Delete(namespace, name)
}
