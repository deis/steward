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

// newCataloger creates a new Tiller-backed mode.Cataloger
func newCataloger(cfg *config) mode.Cataloger {
	return cataloger{
		svc: &mode.Service{
			ServiceInfo: mode.ServiceInfo{
				ID:          cfg.ServiceID,
				Name:        cfg.ServiceName,
				Description: cfg.ServiceDescription,
			},
			Plans: []mode.ServicePlan{
				mode.ServicePlan{
					ID:          cfg.PlanID,
					Name:        cfg.PlanName,
					Description: cfg.PlanDescription,
				},
			},
		},
	}
}
