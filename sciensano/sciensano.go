package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"net/http"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	apiclient.APIClient
}

// APIClient exposes the supported Sciensano APIs
type APIClient interface {
	TestsGetter
	VaccinationGetter
}

var _ APIClient = &Client{}

// NewCachedClient creates a new Client which caches results for duration interval
func NewCachedClient(duration time.Duration) *Client {
	return &Client{
		APIClient: &apiclient.Cache{
			APIClient: &apiclient.Client{HTTPClient: &http.Client{}},
			Retention: duration,
		},
	}
}
