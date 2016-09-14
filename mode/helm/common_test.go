package helm

import (
	"path/filepath"

	"github.com/juju/loggo"
)

func init() {
	logger.SetLogLevel(loggo.TRACE)
}

func alpineChartLoc() string {
	return filepath.Join("..", "..", "testdata", "alpine-0.1.1.tgz")
}
