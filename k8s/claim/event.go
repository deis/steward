package claim

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

// Event represents the event that a service plan claim has changed in kubernetes. It implements fmt.Stringer
type Event struct {
	claim     *ServicePlanClaimWrapper
	operation watch.EventType
}

func eventFromConfigMapEvent(raw watch.Event) (*Event, error) {
	configMap, ok := raw.Object.(*api.ConfigMap)
	if !ok {
		return nil, errNotAConfigMap
	}
	claimWrapper, err := servicePlanClaimWrapperFromConfigMap(configMap)
	if err != nil {
		return nil, err
	}
	return &Event{
		claim:     claimWrapper,
		operation: raw.Type,
	}, nil
}

func (e Event) toConfigMap() *api.ConfigMap {
	return e.claim.toConfigMap()
}

// String is the fmt.Stringer interface implementation
func (e Event) String() string {
	return fmt.Sprintf("%s %s", e.operation, *e.claim)
}
