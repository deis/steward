package mode

// Lifecycler is a composition of the provisioner, deprovisioner, binder and unbinder. It's intended for use in passing to functions that require all functionality
type Lifecycler interface {
	Provisioner
	Deprovisioner
	Binder
	Unbinder
}

// Cataloger lists all the available services
type Cataloger interface {
	List() ([]*Service, error)
}

// Provisioner provisions services, according to the mode implementation
type Provisioner interface {
	Provision(instanceID string, req *ProvisionRequest) (*ProvisionResponse, error)
}

// Deprovisioner deprovisions services, according to the mode implementation
type Deprovisioner interface {
	Deprovision(instanceID, serviceID, planID string) (*DeprovisionResponse, error)
}

// Binder binds services to apps, according to the mode implementation
type Binder interface {
	Bind(instanceID, bindingID string, bindRequest *BindRequest) (*BindResponse, error)
}

// Unbinder unbinds services from apps, according to the mode implementation
type Unbinder interface {
	Unbind(serviceID, planID, instanceID, bindingID string) error
}
