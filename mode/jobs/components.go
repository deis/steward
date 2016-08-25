package jobs

import (
	"github.com/deis/steward/mode"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

// GetComponents returns suitable implementations of the Cataloger and Lifecycler interfaces
func GetComponents(cl *kcl.Client) (mode.Cataloger, *mode.Lifecycler, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, nil, err
	}
	logger.Infof(
		"starting in jobs mode with image %s",
		cfg.Image,
	)
	pr := newPodRunner(cl, cfg)
	cataloger := newCataloger(pr)
	lifecycler := newLifecycler(pr)
	return cataloger, lifecycler, nil
}
