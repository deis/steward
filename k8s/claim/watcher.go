package claim

// Watcher represents the watcher that returns service plan claims and returns them in ResultChan
type Watcher interface {
	ResultChan() <-chan *Event
}
