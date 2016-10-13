package k8s

const (
	// StatusProvisioning is the status indicating that the provisioning process has started
	StatusProvisioning ServicePlanClaimStatus = "provisioning"
	// StatusProvisioned is the status indicating that the provisioning process has succeeded
	StatusProvisioned ServicePlanClaimStatus = "provisioned"
	// StatusBinding is the status indicating that the binding process has started
	StatusBinding ServicePlanClaimStatus = "binding"
	// StatusBound is the status indicating that the binding process has succeeded
	StatusBound ServicePlanClaimStatus = "bound"
	// StatusUnbinding is the status indicating that the unbinding process has started
	StatusUnbinding ServicePlanClaimStatus = "unbinding"
	// StatusUnbound is the status indicating that the unbinding process has succeeded
	StatusUnbound ServicePlanClaimStatus = "unbound"
	// StatusDeprovisioning is the status indicating that the deprovisioning process has started
	StatusDeprovisioning ServicePlanClaimStatus = "deprovisioning"
	// StatusDeprovisioned is the status indicating that the deprovisioning process has succeeded
	StatusDeprovisioned ServicePlanClaimStatus = "deprovisioned"
	// StatusFailed is the status indicating that a service's creation or deletion operation has failed for some reason. The human-readable explanation of the failure will be written to the status description
	StatusFailed ServicePlanClaimStatus = "failed"
)

// ServicePlanClaimStatus is the type representing the current status of a claim. It implements fmt.Stringer
type ServicePlanClaimStatus string

// StringIsStatus returns true if s == st.String()
func StringIsStatus(s string, st ServicePlanClaimStatus) bool {
	return s == st.String()
}

// String is the fmt.Stringer interface implementation
func (s ServicePlanClaimStatus) String() string {
	return string(s)
}
