package reporter

import (
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient/cache"
	sciensanoClient "github.com/clambin/sciensano/apiclient/sciensano"
	vaccinesClient "github.com/clambin/sciensano/apiclient/vaccines"
	cache2 "github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/cases"
	"github.com/clambin/sciensano/reporter/hospitalisations"
	"github.com/clambin/sciensano/reporter/mortality"
	"github.com/clambin/sciensano/reporter/testresults"
	"github.com/clambin/sciensano/reporter/vaccinations"
	"github.com/clambin/sciensano/reporter/vaccines"
	"time"
)

// Client queries different Reporter APIs
type Client struct {
	APICache         cache.Holder
	ReportCache      *cache2.Cache
	Cases            cases.Reporter
	Hospitalisations hospitalisations.Reporter
	Mortality        mortality.Reporter
	TestResults      testresults.Reporter
	Vaccinations     vaccinations.Reporter
	Vaccines         vaccines.Reporter
}

// New creates a new Client which caches results for duration interval
func New(duration time.Duration) *Client {
	return NewWithOptions(duration, client.Options{
		PrometheusMetrics: client.NewMetrics("sciensano", ""),
	})
}

// NewWithOptions creates a new Client with the provided Options
func NewWithOptions(duration time.Duration, options client.Options) *Client {
	apiCache := &cache.Cache{
		Fetchers: []cache.Fetcher{
			&sciensanoClient.Client{
				Caller: &client.InstrumentedClient{
					Options:     options,
					Application: "sciensanoClient",
				},
			},
			&vaccinesClient.Client{
				Caller: &client.InstrumentedClient{
					Options:     options,
					Application: "vaccines",
				},
			},
		},
	}
	reportsCache := cache2.NewCache(duration)
	return &Client{
		APICache:         apiCache,
		ReportCache:      reportsCache,
		Cases:            cases.Reporter{ReportCache: reportsCache, APICache: apiCache},
		Hospitalisations: hospitalisations.Reporter{ReportCache: reportsCache, APICache: apiCache},
		Mortality:        mortality.Reporter{ReportCache: reportsCache, APICache: apiCache},
		TestResults:      testresults.Reporter{ReportCache: reportsCache, APICache: apiCache},
		Vaccinations:     vaccinations.Reporter{ReportCache: reportsCache, APICache: apiCache},
		Vaccines:         vaccines.Reporter{ReportCache: reportsCache, APICache: apiCache},
	}
}
