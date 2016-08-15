package testutil

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/arschles/testsrv"
)

var (
	errMalformedURLString = errors.New("malformed URL string")
)

// HostAndPort returns the host and port represented in srv.URLStr(). returns a non-nil error if the URL string was malformed or it wasn't but the port was not a number
func HostAndPort(srv *testsrv.Server) (string, int, error) {
	u, err := url.Parse(srv.URLStr())
	spl := strings.Split(u.Host, ":")
	if len(spl) != 2 {
		return "", 0, errMalformedURLString
	}
	port, err := strconv.Atoi(spl[1])
	if err != nil {
		return "", 0, err
	}
	return spl[0], port, nil
}
