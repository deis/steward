package claim

import (
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// InteractorNamespacer gets the Interface for a given namespace
type InteractorNamespacer interface {
	Interactor(string) Interactor
}

type cmNamespacer struct {
	cmNS kcl.ConfigMapsNamespacer
}

// NewConfigMapsInteractorNamespacer returns a new EventsNamespacer that works on config maps
func NewConfigMapsInteractorNamespacer(cmns kcl.ConfigMapsNamespacer) InteractorNamespacer {
	return cmNamespacer{cmNS: cmns}
}

func (c cmNamespacer) Interactor(ns string) Interactor {
	return cmInterface{cm: c.cmNS.ConfigMaps(ns)}
}
