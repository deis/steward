package fake

// Unbinder is a fake (github.com/deis/steward/mode).Unbinder implementation, suitable for use in unit tests
type Unbinder struct {
}

// Unbind is the Unbinder interface implementaion. It returns nil
func (u *Unbinder) Unbind(serviceID, planID, instanceID, bindingID string) error {
	return nil
}
