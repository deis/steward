package mode

// Lifecycler is a composition of the provisioner, deprovisioner, binder and unbinder. It's intended for use in passing to functions that require all functionality
type Lifecycler struct {
	Provisioner
	Deprovisioner
	Binder
	Unbinder
}
