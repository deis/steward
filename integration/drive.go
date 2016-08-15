package main

import (
	"fmt"
	"sync"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
	"github.com/pborman/uuid"
)

type multipleErrors struct {
	errs []error
}

func (m multipleErrors) Error() string {
	return fmt.Sprintf("%d error(s): %s", len(m.errs), m.errs)
}

func drive(cl *cf.RESTClient, targetNS, targetName string) error {
	cataloger := cf.NewCataloger(cl)
	provisioner := cf.NewProvisioner(cl)
	deprovisioner := cf.NewDeprovisioner(cl)

	svcs, err := cataloger.List()
	if err != nil {
		return err
	}
	logger.Debugf("found %d service(s)", len(svcs))

	var wg sync.WaitGroup
	errCh := make(chan error)
	doneCh := make(chan struct{})

	for _, svc := range svcs {
		for _, plan := range svc.Plans {
			wg.Add(1)
			s := *svc
			p := plan
			go func(svc *mode.Service, plan *mode.ServicePlan) {
				logger.Debugf("starting service %s, plan %s", svc.Name, plan.Name)
				defer wg.Done()
				instID := uuid.New()
				bindID := uuid.New()

				logger.Debugf("service %s, plan %s provisioning", svc.Name, plan.Name)
				if _, _, err := provision(provisioner, svc.ID, plan.ID, instID); err != nil {
					errCh <- err
					return
				}

				logger.Debugf("service %s, plan %s binding", svc.Name, plan.Name)
				if _, err := bind(cl, svc.ID, plan.ID, instID, bindID, targetNS, targetName); err != nil {
					errCh <- err
					return
				}

				logger.Debugf("service %s, plan %s unbinding", svc.Name, plan.Name)
				if err := unbind(cl, svc.ID, plan.ID, instID, bindID, targetNS, targetName); err != nil {
					errCh <- err
					return
				}

				logger.Debugf("service %s, plan %s deprovisioning", svc.Name, plan.Name)
				if _, err := deprovision(deprovisioner, svc.ID, plan.ID, instID); err != nil {
					errCh <- err
					return
				}

			}(&s, &p)
		}
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	errs := []error{}
	select {
	case err := <-errCh:
		errs = append(errs, err)
	case <-doneCh:
	}
	if len(errs) == 0 {
		return nil
	}
	return multipleErrors{errs: errs}
}
