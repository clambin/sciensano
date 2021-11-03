package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"net/http"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	apiclient.Getter

	cache *Cache
}

// APIClient exposes the supported Sciensano APIs
type APIClient interface {
	TestResultsGetter
	VaccinationGetter
	CasesGetter
	MortalityGetter
}

var _ APIClient = &Client{}

// NewCachedClient creates a new Client which caches results for duration interval
func NewCachedClient(duration time.Duration) *Client {
	return &Client{
		Getter: &apiclient.Cache{
			Getter:    &apiclient.Client{HTTPClient: &http.Client{}},
			Retention: duration,
		},

		cache: NewCache(duration),
	}
}
