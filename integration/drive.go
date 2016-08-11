package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/mode/cf"
	"github.com/juju/loggo"
	"github.com/pborman/uuid"
)

type multipleErrors struct {
	errs []error
}

func (m multipleErrors) Error() string {
	return fmt.Sprintf("%d error(s): %s", len(m.errs), m.errs)
}

func drive(logger loggo.Logger, cl *cf.RESTClient, targetNS string) error {
	cataloger := cf.NewCataloger(logger, cl)
	provisioner := cf.NewProvisioner(logger, cl)
	deprovisioner := cf.NewDeprovisioner(logger, cl)

	svcs, err := cataloger.List()
	if err != nil {
		return err
	}
	log.Printf("found %d service(s)", len(svcs))

	var wg sync.WaitGroup
	errCh := make(chan error)
	doneCh := make(chan struct{})

	for svcIdx, svc := range svcs {
		for planIdx, plan := range svc.Plans {
			wg.Add(1)
			go func(svcIdx, planIdx int, svc *mode.Service, plan *mode.ServicePlan) {
				defer wg.Done()
				instID := uuid.New()
				bindID := uuid.New()

				logger.Debugf("service %d (%s), plan %d (%s) provisioning", svcIdx, svc.Name, planIdx, plan.Name)
				if _, _, err := provision(logger, provisioner, svc.ID, plan.ID, instID); err != nil {
					errCh <- err
					return
				}

				logger.Debugf("service %d (%s), plan %d (%s) binding", svcIdx, svc.Name, planIdx, plan.Name)
				if _, err := bind(logger, cl, svc.ID, plan.ID, instID, bindID, targetNS); err != nil {
					errCh <- err
					return
				}

				logger.Debugf("service %d (%s), plan %d (%s) unbinding", svcIdx, svc.Name, planIdx, plan.Name)
				if err := unbind(logger, cl, svc.ID, plan.ID, instID, bindID, targetNS); err != nil {
					errCh <- err
					return
				}

				logger.Debugf("service %d (%s), plan %d (%s) deprovisioning", svcIdx, svc.Name, planIdx, plan.Name)
				if _, err := deprovision(logger, deprovisioner, svc.ID, plan.ID, instID); err != nil {
					errCh <- err
					return
				}

			}(svcIdx, planIdx, svc, &plan)
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
