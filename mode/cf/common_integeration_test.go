// +build integration

package cf

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/deis/steward/test-utils/k8s"
	testsetup "github.com/deis/steward/test-utils/setup"
	"github.com/technosophos/moniker"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/util/intstr"
	"k8s.io/client-go/1.4/pkg/util/wait"
)

var namespace string

func TestMain(m *testing.M) {
	testsetup.SetupAndTearDown(m, setup, teardown)
}

func setup() error {
	namespace = moniker.New().NameSep("-")
	if err := k8s.EnsureNamespace(namespace); err != nil {
		return err
	}
	if err := ensureBroker(); err != nil {
		return err
	}
	return nil
}

// ensureBroker sets up a CF broker within the leased cluster
func ensureBroker() error {
	clientset, err := k8s.GetClientset()
	if err != nil {
		return err
	}
	serviceClient := clientset.Services(namespace)
	if _, err = serviceClient.Create(&v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      "cf-sample-broker",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "cf-sample-broker",
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				},
			},
			Selector: map[string]string{
				"app": "cf-sample-broker",
			},
		},
	}); err != nil {
		return err
	}
	podClient := clientset.Pods(namespace)
	if _, err = podClient.Create(&v1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "cf-sample-broker",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "cf-sample-broker",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Name:            "cf-sample-broker",
					Image:           "quay.io/deisci/cf-sample-broker:devel",
					ImagePullPolicy: v1.PullAlways,
					LivenessProbe: &v1.Probe{
						FailureThreshold:    3,
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
						SuccessThreshold:    1,
						Handler: v1.Handler{
							TCPSocket: &v1.TCPSocketAction{
								Port: intstr.FromInt(8080),
							},
						},
					},
					ReadinessProbe: &v1.Probe{
						FailureThreshold:    1,
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
						SuccessThreshold:    1,
						Handler: v1.Handler{
							TCPSocket: &v1.TCPSocketAction{
								Port: intstr.FromInt(8080),
							},
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	// Wait for service to get an external IP
	wg.Add(1)
	var service *v1.Service
	go func() {
		defer wg.Done()
		if err := wait.PollImmediate(time.Duration(2)*time.Second, time.Duration(5)*time.Minute, func() (bool, error) {
			service, err = serviceClient.Get("cf-sample-broker")
			if err != nil {
				return false, err
			}
			return len(service.Status.LoadBalancer.Ingress) > 0, nil
		}); err != nil {
			errCh <- err
		}
	}()
	// Wait for CF broker pod to be running and ready
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := wait.PollImmediate(time.Duration(2)*time.Second, time.Duration(5)*time.Minute, func() (bool, error) {
			pod, err := podClient.Get("cf-sample-broker")
			if err != nil {
				return false, err
			}
			if pod.Status.Phase == v1.PodRunning {
				var ready bool
				for _, condition := range pod.Status.Conditions {
					if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
						ready = true
						break
					}
				}
				return ready, nil
			}
			return false, nil
		}); err != nil {
			errCh <- err
		}
	}()
	// Wait for the service to have endpoints
	endpointsClient := clientset.Endpoints(namespace)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := wait.PollImmediate(time.Duration(2)*time.Second, time.Duration(5)*time.Minute, func() (bool, error) {
			endpoints, err := endpointsClient.Get("cf-sample-broker")
			if err != nil {
				log.Println(".")
				return false, err
			}
			return len(endpoints.Subsets) > 0 && len(endpoints.Subsets[0].Addresses) > 0, nil
		}); err != nil {
			errCh <- err
		}
	}()
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		close(errCh)
	case err := <-errCh:
		close(errCh)
		return err
	}
	// krancour: I wish we didn't have to do this, but even with all the polling done above to
	// ensure all pieces of the sample broker are ready, I've observed case after case where the
	// initial request to the sample broker takes many seconds (which triggers timeouts in some test
	// cases). This would seem to indicates that SOMETHING is indeed NOT ready after all. So,
	// reluctantly, we're waiting an extra 30 seconds here. With This extra bit of padding, the
	// first response from the sample broker is reliably received in little more than 100 ms.
	time.Sleep(time.Duration(30) * time.Second)
	os.Setenv("STEWARD_CF_BROKER_SCHEME", "http")
	os.Setenv("CF_BROKER_HOSTNAME", service.Status.LoadBalancer.Ingress[0].IP)
	os.Setenv("CF_BROKER_PORT", "80")
	os.Setenv("CF_BROKER_USERNAME", "admin")
	os.Setenv("CF_BROKER_PASSWORD", "password")
	return nil
}

func teardown() error {
	// This will also delete the broker
	if err := k8s.DeleteNamespace(namespace); err != nil {
		return err
	}
	return nil
}
