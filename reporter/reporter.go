package reporter

import (
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/measurement"
	"net/http"
	"time"
)

// Client queries different Reporter APIs
type Client struct {
	APICache    measurement.Holder
	ReportCache *Cache
}

// Reporter exposes the supported Reporter APIs
type Reporter interface {
	TestResultsGetter
	VaccinationGetter
	CasesGetter
	MortalityGetter
	HospitalisationsGetter
	VaccinesGetter
}

var _ Reporter = &Client{}

// New creates a new Client which caches results for duration interval
func New(duration time.Duration) *Client {
	return &Client{
		APICache: &measurement.Cache{
			Fetchers: []measurement.Fetcher{
				&sciensano.Client{HTTPClient: &http.Client{}},
				&vaccines.Client{HTTPClient: &http.Client{}},
			},
		},
		ReportCache: NewCache(duration),
	}
}
