package cmd

import (
	"github.com/deis/steward/mode"
	"k8s.io/client-go/1.4/kubernetes"
)

// GetComponents returns suitable implementations of the Cataloger and Lifecycler interfaces
func GetComponents(cl *kubernetes.Clientset) (mode.Cataloger, *mode.Lifecycler, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, nil, err
	}
	logger.Infof(
		"starting in cmd mode with image %s",
		cfg.Image,
	)
	pr := newPodRunner(cl, cfg)
	cataloger := newCataloger(pr)
	lifecycler := newLifecycler(pr)
	return cataloger, lifecycler, nil
}
