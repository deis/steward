package claim

import (
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
)

// InteractorNamespacer gets the Interface for a given namespace
type InteractorNamespacer interface {
	Interactor(string) Interactor
}

type cmNamespacer struct {
	cmNS v1.ConfigMapsGetter
}

// NewConfigMapsInteractorNamespacer returns a new EventsNamespacer that works on config maps
func NewConfigMapsInteractorNamespacer(cmns v1.ConfigMapsGetter) InteractorNamespacer {
	return cmNamespacer{cmNS: cmns}
}

func (c cmNamespacer) Interactor(ns string) Interactor {
	return cmInterface{cm: c.cmNS.ConfigMaps(ns)}
}
