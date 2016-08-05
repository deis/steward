package mode

// Service is the represntation of a steward service. It also is compatible with the CloudFoundry catalog API. See https://docs.cloudfoundry.org/services/api.html#catalog-mgmt for more detail
type Service struct {
	ServiceInfo
	Plans []ServicePlan `json:"plans"`
}
