package cloudlycke

import (
	"net/http"
)

func newCloudlyckeClient() *http.Client {
	return &http.Client{}
}