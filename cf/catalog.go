package cf

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juju/loggo"
)

type errUnexpectedResponseCode struct {
	url      string
	actual   int
	expected int
}

func (e errUnexpectedResponseCode) Error() string {
	return fmt.Sprintf("%s - expected response code %d, actual %d", e.url, e.expected, e.actual)
}

// GetCatalog fetches the catalog of services from baseURL, which is expected to be a valid CloudFoundry service broker API. It assumes that the API version is 2.9. See https://docs.cloudfoundry.org/services/api.html for more detail
func GetCatalog(logger loggo.Logger, cl *Client) ([]Service, error) {
	req, err := cl.Get(logger, "v2", "catalog")
	if err != nil {
		return nil, err
	}
	res, err := cl.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	serviceList := new(serviceList)
	// TODO: drain the response body to avoid a connection leak
	if err := json.NewDecoder(res.Body).Decode(serviceList); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errUnexpectedResponseCode{expected: http.StatusOK, actual: res.StatusCode}
	}
	return serviceList.Services, nil
}
