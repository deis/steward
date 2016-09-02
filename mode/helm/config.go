package helm

import (
	"github.com/deis/steward/config"
)

// Config is the envconfig-compatible struct for a backing Tiller server
type Config struct {
	TillerIP           string `envconfig:"HELM_TILLER_IP" required:"true"`
	TillerPort         int    `envconfig:"HELM_TILLER_PORT" required:"true"`
	ChartURL           string `envconfig:"HELM_CHART_URL" required:"true"`
	ChartInstallNS     string `envconfig:"HELM_CHART_INSTALL_NAMESPACE" required:"true"`
	ProvisionBehavior  string `envconfig:"HELM_PROVISION_BEHAVIOR" required:"true"`
	ServiceID          string `envconfig:"HELM_SERVICE_ID" required:"true"`
	ServiceName        string `envconfig:"HELM_SERVICE_NAME" required:"true"`
	ServiceDescription string `envconfig:"HELM_SERVICE_DESCRIPTION" required:"true"`
	PlanID             string `envconfig:"HELM_PLAN_ID" required:"true"`
	PlanName           string `envconfig:"HELM_PLAN_NAME" required:"true"`
	PlanDescription    string `envconfig:"HELM_PLAN_DESCRIPTION" required:"true"`
}

// GetConfig gets the configuration for helm mode
func GetConfig() (*Config, error) {
	ret := new(Config)
	if err := config.Load(ret); err != nil {
		return nil, err
	}
	return ret, nil
}
