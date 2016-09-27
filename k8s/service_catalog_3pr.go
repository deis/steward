package k8s

import (
	"k8s.io/client-go/1.4/pkg/api/unversioned"
	"k8s.io/client-go/1.4/pkg/api/v1"
	ext "k8s.io/client-go/1.4/pkg/apis/extensions/v1beta1"
)

// ServiceCatalog3PR is the struct representation of a Third Party Resource
var ServiceCatalog3PR = &ext.ThirdPartyResource{
	TypeMeta: unversioned.TypeMeta{
		Kind:       "ThirdPartyResource",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "service-catalog-entry.steward.deis.io",
		Labels: map[string]string{
			"heritage": "deis",
		},
	},
	Description: "A description of a single (service, plan) pair that a steward instance is able to provision",
	Versions: []ext.APIVersion{
		{Name: "v1"},
	},
}
