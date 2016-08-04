package broker

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/deis/steward/k8s"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	configMapKind = "ConfigMap"
	secretKind    = "Secret"
)

func getObjectMeta(namespace, serviceID, planID, kind string) api.ObjectMeta {
	standardLabels := map[string]string{
		"created-by": "steward",
		"created-at": time.Now().String(),
	}
	return api.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s-%s", serviceID, planID, kind),
		Namespace: namespace,
		Labels:    standardLabels,
	}
}

// writes everything in creds to a configMap and secret. returns the fully qualified name of the configMap and secrets (respectively)
func writeToKubernetes(
	serviceID,
	planID,
	namespace string,
	creds bindingCredentials,
	configMapCreator k8s.ConfigMapCreator,
	secretCreator k8s.SecretCreator,
) (*qualifiedName, []*qualifiedName, error) {
	configMap := &api.ConfigMap{
		TypeMeta:   unversioned.TypeMeta{Kind: configMapKind},
		ObjectMeta: getObjectMeta(namespace, serviceID, planID, configMapKind),
		Data: map[string]string{
			"uri":      creds.URI,
			"hostname": creds.Hostname,
			"port":     strconv.Itoa(creds.Port),
			"name":     creds.Name,
			"username": creds.Username,
		},
	}
	if _, err := configMapCreator(namespace, configMap); err != nil {
		return nil, nil, err
	}
	configMapQualifiedName := &qualifiedName{
		Name:      configMap.ObjectMeta.Name,
		Namespace: namespace,
	}

	encodedPassword := []byte(base64.StdEncoding.EncodeToString([]byte(creds.Password)))
	secret := &api.Secret{
		TypeMeta:   unversioned.TypeMeta{Kind: secretKind},
		ObjectMeta: getObjectMeta(namespace, serviceID, planID, secretKind),
		Type:       api.SecretTypeOpaque,
		Data:       map[string][]byte{"password": encodedPassword},
	}
	if _, err := secretCreator(namespace, secret); err != nil {
		return nil, nil, err
	}
	secretQualifiedNames := []*qualifiedName{
		&qualifiedName{Name: configMap.ObjectMeta.Name, Namespace: namespace},
	}

	return configMapQualifiedName, secretQualifiedNames, nil
}
