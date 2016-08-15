package testutil

import (
	"github.com/juju/loggo"
)

// ConfigLogger configures the root logger to be verbose and suitable for tests.
func ConfigLogger() {
	loggo.GetLogger("").SetLogLevel(loggo.TRACE)
}
