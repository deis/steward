package mode

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
	//TODO
}

// Binder binds services to apps, according to the mode implementation
type Binder interface {
	Bind(instanceID, bindingID string, bindRequest *BindRequest) (*BindResponse, error)
}

// Unbinder unbinds services from apps, according to the mode implementation
type Unbinder interface {
	//TODO
}
