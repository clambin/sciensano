package reporter

import (
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient/cache"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"time"
)

// Client queries different Reporter APIs
type Client struct {
	APICache    cache.Holder
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
	return NewWithOptions(duration, client.Options{
		PrometheusMetrics: client.NewMetrics("sciensano", ""),
	})
}

// NewWithOptions creates a new Client with the provided Options
func NewWithOptions(duration time.Duration, options client.Options) *Client {
	return &Client{
		APICache: &cache.Cache{
			Fetchers: []cache.Fetcher{
				&sciensano.Client{
					Caller: &client.InstrumentedClient{
						Options:     options,
						Application: "sciensano",
					},
				},
				&vaccines.Client{
					Caller: &client.InstrumentedClient{
						Options:     options,
						Application: "vaccines",
					},
				},
			},
		},
		ReportCache: NewCache(duration),
	}
}
