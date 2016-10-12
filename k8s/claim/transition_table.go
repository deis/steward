package claim

import (
	"github.com/deis/steward/k8s"
	"github.com/deis/steward/k8s/claim/state"
	"k8s.io/client-go/1.4/pkg/watch"
)

var (
	transitionTable = map[state.Current]nextFunc{
		// if event was ADDED and action is provison, next action is processProvision
		state.NewCurrentNoStatus(k8s.ActionProvision, watch.Added): nextFunc(processProvision),
		// if event was ADDED and action is bind, next action is processBind
		state.NewCurrentNoStatus(k8s.ActionBind, watch.Added): nextFunc(processBind),
		// if event was ADDED and action is create, next action is processProvision
		state.NewCurrentNoStatus(k8s.ActionCreate, watch.Added): nextFunc(processProvision),
		// if event was MODIFIED, status is provisioned and action is create, next action is processBind
		state.NewCurrent(k8s.StatusProvisioned, k8s.ActionCreate, watch.Modified): nextFunc(processBind),
		// if event was MODIFIED, status is provisioned and action is deprovision, next action is processDeprovision
		state.NewCurrent(k8s.StatusProvisioned, k8s.ActionDeprovision, watch.Modified): nextFunc(processDeprovision),
		// if event was MODIFIED, status is provisioned and action is bind, next action is processBind
		state.NewCurrent(k8s.StatusProvisioned, k8s.ActionBind, watch.Modified): nextFunc(processBind),
		// if event was MODIFIED, status is bound and action is delete, next action is processUnbind
		state.NewCurrent(k8s.StatusBound, k8s.ActionDelete, watch.Modified): nextFunc(processUnbind),
		// if event was MODIFIED, status is unbound and action is delete, next action is processDeprovision
		state.NewCurrent(k8s.StatusUnbound, k8s.ActionDelete, watch.Modified): nextFunc(processDeprovision),
		// if event was MODIFIED, status is bound and action is unbind, next action is processUnbind
		state.NewCurrent(k8s.StatusBound, k8s.ActionUnbind, watch.Modified): nextFunc(processUnbind),
		// if event was MODIFIED, status is unbound and action is deprovision, next action is processDeprovision
		state.NewCurrent(k8s.StatusUnbound, k8s.ActionDeprovision, watch.Modified): nextFunc(processDeprovision),
		// if event was MODIFIED, status is unbound and action is bind, next action is processBind
		state.NewCurrent(k8s.StatusUnbound, k8s.ActionBind, watch.Modified): nextFunc(processBind),
	}
)
