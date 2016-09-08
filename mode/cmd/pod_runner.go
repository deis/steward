package cmd

import (
	"fmt"
	"strings"

	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/wait"
)

type podRunner struct {
	client *kcl.Client
	config *config
}

// run runs the provided command with the provided arguments as a k8s pod using this podRunner's
// configured image, configMap, and secret.
func (r *podRunner) run(podName string, command string, args ...string) (string, error) {
	// Get the manifest
	pod, err := r.getManifest(podName, command, args)
	if err != nil {
		return "", err
	}

	podClient := r.client.Pods(r.config.PodNamespace)

	// Schedule the pod
	if _, err := podClient.Create(pod); err != nil {
		return "", err
	}
	defer podClient.Delete(podName, &api.DeleteOptions{})

	// Wait for the pod to be in state succeeded or failed
	if err := r.waitForPodEnd(podName); err != nil {
		return "", fmt.Errorf("error waiting for pod to complete: %s", err)
	}

	// Get the latest pod state
	pod, err = podClient.Get(podName)
	if err != nil {
		return "", fmt.Errorf("error getting pod: %s", err)
	}

	// Get the exit code
	logger.Debugf("Checking pod exit code")
	containerExitCode := pod.Status.ContainerStatuses[0].State.Terminated.ExitCode

	logReq := podClient.GetLogs(podName, &api.PodLogOptions{})
	logRes := logReq.Do()
	logBytes, err := logRes.Raw()
	if err != nil {
		return "", fmt.Errorf("error retrieving pod output: %s", err)
	}
	output := string(logBytes)

	// If the exit code is not 0 return an error.
	if containerExitCode != 0 {
		return output, fmt.Errorf("Pod exited with code %d", containerExitCode)
	}

	return output, nil
}

func (r *podRunner) getManifest(podName string, command string, args []string) (*api.Pod, error) {
	commandAndArgs := []string{command}
	commandAndArgs = append(commandAndArgs, args...)
	envVars, err := r.getEnvManifestFragment()
	if err != nil {
		return nil, err
	}
	volumeMounts := r.getVolumeMountsFragment()
	volumes := r.getVolumesFragment()
	pod := &api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name:      podName,
			Namespace: "steward",
			Labels: map[string]string{
				"heritage": "steward",
			},
		},
		Spec: api.PodSpec{
			RestartPolicy: api.RestartPolicyNever,
			Containers: []api.Container{
				api.Container{
					Name:            podName,
					Image:           r.config.Image,
					ImagePullPolicy: api.PullPolicy("Always"),
					Command:         commandAndArgs,
					Env:             envVars,
					VolumeMounts:    volumeMounts,
				},
			},
			Volumes: volumes,
		},
	}
	return pod, nil
}

func (r *podRunner) getEnvManifestFragment() ([]api.EnvVar, error) {
	configMapKeys, err := r.getKeysFromConfigMap()
	if err != nil {
		return nil, err
	}
	secretKeys, err := r.getKeysFromSecret()
	if err != nil {
		return nil, err
	}
	envVars := make([]api.EnvVar, len(configMapKeys)+len(secretKeys))
	i := 0
	for _, configMapKey := range configMapKeys {
		envVars[i] = api.EnvVar{
			Name: toEnvVarName(configMapKey),
			ValueFrom: &api.EnvVarSource{
				ConfigMapKeyRef: &api.ConfigMapKeySelector{
					LocalObjectReference: api.LocalObjectReference{
						Name: r.config.ConfigMapName,
					},
					Key: configMapKey,
				},
			},
		}
		i++
	}
	for _, secretKey := range secretKeys {
		envVars[i] = api.EnvVar{
			Name: toEnvVarName(secretKey),
			ValueFrom: &api.EnvVarSource{
				SecretKeyRef: &api.SecretKeySelector{
					LocalObjectReference: api.LocalObjectReference{
						Name: r.config.SecretName,
					},
					Key: secretKey,
				},
			},
		}
		i++
	}
	return envVars, nil
}

func toEnvVarName(keyName string) string {
	return strings.Replace(strings.ToUpper(keyName), ".", "_", -1)
}

func (r *podRunner) getKeysFromConfigMap() ([]string, error) {
	if r.config.ConfigMapName == "" {
		return []string{}, nil
	}
	cfCl := r.client.ConfigMaps(r.config.PodNamespace)
	configMap, err := cfCl.Get(r.config.ConfigMapName)
	if err != nil {
		return nil, err
	}
	data := configMap.Data
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (r *podRunner) getKeysFromSecret() ([]string, error) {
	if r.config.SecretName == "" {
		return []string{}, nil
	}
	cfCl := r.client.Secrets(r.config.PodNamespace)
	secret, err := cfCl.Get(r.config.SecretName)
	if err != nil {
		return nil, err
	}
	data := secret.Data
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (r *podRunner) getVolumeMountsFragment() []api.VolumeMount {
	var mounts []api.VolumeMount
	if r.config.ConfigMapName != "" {
		mounts = append(mounts, api.VolumeMount{
			Name:      "config-volume",
			MountPath: "/config",
		})
	}
	if r.config.SecretName != "" {
		mounts = append(mounts, api.VolumeMount{
			Name:      "secret-volume",
			MountPath: "/secret",
		})
	}
	return mounts
}

func (r *podRunner) getVolumesFragment() []api.Volume {
	var volumes []api.Volume
	if r.config.ConfigMapName != "" {
		volumes = append(volumes, api.Volume{
			Name: "config-volume",
			VolumeSource: api.VolumeSource{
				ConfigMap: &api.ConfigMapVolumeSource{
					LocalObjectReference: api.LocalObjectReference{
						Name: r.config.ConfigMapName,
					},
				},
			},
		})
	}
	if r.config.SecretName != "" {
		volumes = append(volumes, api.Volume{
			Name: "secret-volume",
			VolumeSource: api.VolumeSource{
				Secret: &api.SecretVolumeSource{
					SecretName: r.config.SecretName,
				},
			},
		})
	}
	return volumes
}

// waitForPodEnd waits for a pod in state succeeded or failed
func (r *podRunner) waitForPodEnd(podName string) error {
	logger.Debugf(
		"Waiting for pod %s/%s to exit. Checking every %s for %s",
		r.config.PodNamespace,
		podName,
		r.config.getPollInterval(),
		r.config.getTimeout(),
	)
	return wait.PollImmediate(r.config.getPollInterval(), r.config.getTimeout(), func() (bool, error) {
		pod, err := r.client.Pods(r.config.PodNamespace).Get(podName)
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == api.PodSucceeded || pod.Status.Phase == api.PodFailed {
			logger.Debugf("Pod %s/%s has exited", r.config.PodNamespace, podName)
			return true, nil
		}
		return false, nil
	})
}

// newPodRunner builds and returns a podRunner
func newPodRunner(client *kcl.Client, cfg *config) *podRunner {
	return &podRunner{
		client: client,
		config: cfg,
	}
}
