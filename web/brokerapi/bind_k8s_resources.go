package brokerapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	configMapKind = "ConfigMap"
	secretKind    = "Secret"
)

func getResourceName(serviceID, planID, bindingID, instanceID string) string {
	return fmt.Sprintf("%s-%s-%s-%s", serviceID, planID, bindingID, instanceID)
}

func getObjectMeta(namespace, serviceID, planID, bindingID, instanceID string) api.ObjectMeta {
	standardLabels := map[string]string{
		"created-by": "steward",
		"created-at": time.Now().String(),
	}
	return api.ObjectMeta{
		Name:      getResourceName(serviceID, planID, bindingID, instanceID),
		Namespace: namespace,
		Labels:    standardLabels,
	}
}

// writes everything in creds to a configMap and secret. returns the fully qualified name of the configMap and secrets (respectively)
func writeToKubernetes(
	serviceID,
	planID,
	bindingID,
	instanceID,
	namespace string,
	publicCreds mode.JSONObject,
	privateCreds mode.JSONObject,
	configMapCreator k8s.ConfigMapCreator,
	secretCreator k8s.SecretCreator,
) (*qualifiedName, []*qualifiedName, error) {

	configMap := &api.ConfigMap{
		TypeMeta:   unversioned.TypeMeta{Kind: configMapKind},
		ObjectMeta: getObjectMeta(namespace, serviceID, planID, bindingID, instanceID),
		Data:       publicCreds,
	}
	if _, err := configMapCreator.Create(namespace, configMap); err != nil {
		return nil, nil, err
	}
	configMapQualifiedName := &qualifiedName{
		Name:      configMap.ObjectMeta.Name,
		Namespace: namespace,
	}

	privateCredsBytes, err := json.Marshal(privateCreds)
	if err != nil {
		return nil, nil, err
	}
	encodedPrivateCreds := base64.StdEncoding.EncodeToString(privateCredsBytes)
	secret := &api.Secret{
		TypeMeta:   unversioned.TypeMeta{Kind: secretKind},
		ObjectMeta: getObjectMeta(namespace, serviceID, planID, bindingID, instanceID),
		Type:       api.SecretTypeOpaque,
		Data:       map[string][]byte{"password": []byte(encodedPrivateCreds)},
	}
	if _, err := secretCreator.Create(namespace, secret); err != nil {
		return nil, nil, err
	}
	secretQualifiedNames := []*qualifiedName{
		&qualifiedName{Name: configMap.ObjectMeta.Name, Namespace: namespace},
	}

	return configMapQualifiedName, secretQualifiedNames, nil
}

func deleteFromKubernetes(
	serviceID,
	planID,
	bindingID,
	instanceID,
	namespace string,
	configMapDeleter k8s.ConfigMapDeleter,
	secretDeleter k8s.SecretDeleter,
) error {
	if err := configMapDeleter.Delete(
		namespace,
		getResourceName(serviceID, planID, bindingID, instanceID),
	); err != nil {
		return err
	}
	if err := secretDeleter.Delete(
		namespace,
		getResourceName(serviceID, planID, bindingID, instanceID),
	); err != nil {
		return err
	}
	return nil
}
