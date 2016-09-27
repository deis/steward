package claim

import (
	"github.com/deis/steward/k8s/claim/state"
	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/pkg/watch"
)

var (
	transitionTable = map[state.Current]nextFunc{
		// if event was ADDED and action is provison, next action is processProvision
		state.NewCurrentNoStatus(mode.ActionProvision, watch.Added): nextFunc(processProvision),
		// if event was ADDED and action is bind, next action is processBind
		state.NewCurrentNoStatus(mode.ActionBind, watch.Added): nextFunc(processBind),
		// if event was ADDED and action is create, next action is processProvision
		state.NewCurrentNoStatus(mode.ActionCreate, watch.Added): nextFunc(processProvision),
		// if event was MODIFIED, status is provisioned and action is create, next action is processBind
		state.NewCurrent(mode.StatusProvisioned, mode.ActionCreate, watch.Modified): nextFunc(processBind),
		// if event was MODIFIED, status is provisioned and action is deprovision, next action is processDeprovision
		state.NewCurrent(mode.StatusProvisioned, mode.ActionDeprovision, watch.Modified): nextFunc(processDeprovision),
		// if event was MODIFIED, status is provisioned and action is bind, next action is processBind
		state.NewCurrent(mode.StatusProvisioned, mode.ActionBind, watch.Modified): nextFunc(processBind),
		// if event was MODIFIED, status is bound and action is delete, next action is processUnbind
		state.NewCurrent(mode.StatusBound, mode.ActionDelete, watch.Modified): nextFunc(processUnbind),
		// if event was MODIFIED, status is unbound and action is delete, next action is processDeprovision
		state.NewCurrent(mode.StatusUnbound, mode.ActionDelete, watch.Modified): nextFunc(processDeprovision),
		// if event was MODIFIED, status is bound and action is unbind, next action is processUnbind
		state.NewCurrent(mode.StatusBound, mode.ActionUnbind, watch.Modified): nextFunc(processUnbind),
		// if event was MODIFIED, status is unbound and action is deprovision, next action is processDeprovision
		state.NewCurrent(mode.StatusUnbound, mode.ActionDeprovision, watch.Modified): nextFunc(processDeprovision),
		// if event was MODIFIED, status is unbound and action is bind, next action is processBind
		state.NewCurrent(mode.StatusUnbound, mode.ActionBind, watch.Modified): nextFunc(processBind),
	}
)
