package setup

import (
	"log"
	"os"
	"testing"
)

// SetupAndTearDown coordinates test setup, test execution, and teardown. This is typically
// invoked within a package's TestMain function.
func SetupAndTearDown(m *testing.M, setup func() error, teardown func() error) {
	var exitCode int
	if err := setup(); err != nil {
		// This isn't immediately fatal, because we still want we want to try and teardown whatever
		// setup did succeed
		log.Printf("cf package encountered fatal error during test setup: %s", err)
		exitCode = 1
	} else {
		exitCode = m.Run()
	}
	// This isn't immediately fatal, because a frequent cleanup step in teardown is deleting a k8s
	// namespace, and that seems prone to spurious errors, even though the operation always seems
	// to succeed. Assuming the tests all passed, we do not want a minor issue such as this to force
	// the test suite to exit non-zero.
	if err := teardown(); err != nil {
		log.Printf("WARNING: cf package encountered error during test teardown: %s", err)
	}
	os.Exit(exitCode)
}
