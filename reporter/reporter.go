package reporter

import (
	"github.com/clambin/go-metrics/caller"
	"github.com/clambin/sciensano/apiclient/cache"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/metrics"
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
	return &Client{
		APICache: &cache.Cache{
			Fetchers: []cache.Fetcher{
				&sciensano.Client{
					Caller: &caller.InstrumentedClient{
						Options:     caller.Options{PrometheusMetrics: metrics.Metrics},
						Application: "sciensano",
					},
				},
				&vaccines.Client{
					Caller: &caller.InstrumentedClient{
						Options:     caller.Options{PrometheusMetrics: metrics.Metrics},
						Application: "vaccines",
					},
				},
			},
		},
		ReportCache: NewCache(duration),
	}
}
