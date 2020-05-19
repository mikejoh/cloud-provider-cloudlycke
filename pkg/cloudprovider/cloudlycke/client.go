package cloudlycke

import (
	"net/http"
)

// newCloudlyckeClient returns a specific HTTP client used when communicating with the Cloudlycke API(s)
func newCloudlyckeClient() *http.Client {
	return &http.Client{}
}
