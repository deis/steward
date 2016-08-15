package testutil

import (
	"github.com/juju/loggo"
)

// GetLogger returns a verbose logger, suitable for tests
func GetLogger() loggo.Logger {
	l := loggo.GetLogger("test")
	l.SetLogLevel(loggo.TRACE)
	return l
}
