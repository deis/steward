package fake

// Lifecycler is a fake implementation of a (github.com/deis/steward/mode).Lifecycler, suitable for use in unit tests
type Lifecycler struct {
	*Provisioner
	*Binder
	*Unbinder
	*Deprovisioner
}
