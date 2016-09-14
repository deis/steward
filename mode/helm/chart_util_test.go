package helm

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/arschles/testsrv"
)

var (
	bgCtx = context.Background()
)

func stdHandler(b []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, bytes.NewBuffer(b))
	})
}

func notFoundHandler([]byte) http.Handler {
	return http.NotFoundHandler()
}

func newChartServer(createHdl func([]byte) http.Handler) (*testsrv.Server, error) {
	fb, err := ioutil.ReadFile(alpineChartLoc())
	if err != nil {
		return nil, err
	}

	hdl := createHdl(fb)
	return testsrv.StartServer(hdl), nil
}

func TestGetChartCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(bgCtx)
	cancel()
	srv, err := newChartServer(stdHandler)
	assert.NoErr(t, err)
	defer srv.Close()
	_, tmpDir, err := getChart(ctx, http.DefaultClient, srv.URLStr())
	defer os.RemoveAll(tmpDir)
	assert.True(t, err != nil, "gettting chart with cancelled context returned no error")
	srv.AcceptN(1, 1*time.Second)
}

func TestGetChartNotFound(t *testing.T) {
	srv, err := newChartServer(notFoundHandler)
	assert.NoErr(t, err)
	defer srv.Close()
	_, _, err = getChart(bgCtx, http.DefaultClient, srv.URLStr())
	if _, ok := err.(errChartNotFound); !ok {
		t.Fatalf("returned error was not a errChartNotFound")
	}
}

func TestGetChart(t *testing.T) {
	srv, err := newChartServer(stdHandler)
	assert.NoErr(t, err)
	defer srv.Close()
	go func() {
		srv.AcceptN(1, 1*time.Second)
	}()
	ch, tmpDir, err := getChart(bgCtx, http.DefaultClient, srv.URLStr())
	defer os.RemoveAll(tmpDir)
	assert.NoErr(t, err)
	assert.NotNil(t, ch, "chart")
}
