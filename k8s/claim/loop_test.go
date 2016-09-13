package claim

import (
	"context"
	"testing"

	"github.com/deis/steward/k8s"
	"github.com/deis/steward/mode"
)

func TestStartLoop(t *testing.T) {
	t.Skip("TODO")
}

func TestReceiveEvent(t *testing.T) {
	ctx := context.Background()
	evt := getEvent(getClaim(mode.ActionProvision))
	iface := &FakeInteractor{}
	secretsNamespacer := &k8s.FakeSecretsNamespacer{}
	lookup := k8s.NewServiceCatalogLookup(nil) // TODO: add service/plan to the catalog
	lifecycler := &mode.Lifecycler{}
	receiveEvent(ctx, evt, iface, secretsNamespacer, lookup, lifecycler)
}

func TestStopLoop(t *testing.T) {
	t.Skip("TODO")
}

func TestWatchChanClosed(t *testing.T) {
	t.Skip("TODO")
}
