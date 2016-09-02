package helm

import (
	"github.com/deis/steward/mode"
)

type cataloger struct {
	svc *mode.Service
}

func (c cataloger) List() ([]*mode.Service, error) {
	return []*mode.Service{c.svc}, nil
}

// NewCataloger creates a new Tiller-backed mode.Cataloger
func NewCataloger(serviceID, serviceName, serviceDescription, planID, planName, planDescription string) mode.Cataloger {
	return cataloger{
		svc: &mode.Service{
			ServiceInfo: mode.ServiceInfo{
				ID:          serviceID,
				Name:        serviceName,
				Description: serviceDescription,
			},
			Plans: []mode.ServicePlan{
				mode.ServicePlan{
					ID:          planID,
					Name:        planName,
					Description: planDescription,
				},
			},
		},
	}
}
