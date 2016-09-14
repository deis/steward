package helm

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deis/steward/web/ctxhttp"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type errChartNotFound struct {
	chartURL string
	code     int
}

func (e errChartNotFound) Error() string {
	return fmt.Sprintf("chart at %s not found (status code %d)", e.chartURL, e.code)
}

// getChart downloads the chart at chartURL to a directory, parses it into a *chart.Chart and returns it along with the root directory of the directory the chart was downloaded to. Returns a non-nil error if the parsing failed. It's the caller's responsibility to delete the chart directory when done with it.
func getChart(ctx context.Context, httpCl *http.Client, chartURL string) (*chart.Chart, string, error) {
	logger.Debugf("downloading chart from %s", chartURL)
	resp, err := ctxhttp.Get(ctx, httpCl, chartURL)
	if err != nil {
		logger.Errorf("downloading chart from %s (%s)", chartURL, err)
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("got status code %d trying to download chart at %s", resp.StatusCode, chartURL)
		return nil, "", errChartNotFound{chartURL: chartURL, code: resp.StatusCode}
	}
	tmpDir, err := ioutil.TempDir("", chartTmpDirPrefix)
	if err != nil {
		return nil, "", err
	}

	fullChartPath := filepath.Join(tmpDir, tmpChartName)
	fd, err := os.Create(fullChartPath)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if err := fd.Close(); err != nil {
			logger.Errorf("closing file descriptor for chart at %s (%s)", fullChartPath, err)
		}
	}()
	logger.Debugf("copying chart to %s", fullChartPath)
	if _, err := io.Copy(fd, resp.Body); err != nil {
		logger.Errorf("copying chart contents to %s (%s)", fullChartPath, err)
		return nil, "", err
	}
	logger.Debugf("loading chart from %s on disk", fullChartPath)
	chart, err := chartutil.Load(fullChartPath)
	if err != nil {
		logger.Errorf("loading chart from %s on disk (%s)", fullChartPath, err)
		return nil, "", err
	}
	return chart, tmpDir, nil
}
