package reporter

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient/fetcher"
	sciensanoClient "github.com/clambin/sciensano/apiclient/sciensano"
	vaccinesClient "github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/reporter/cache"
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
	ReportCache      *cache.Cache
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
func NewWithOptions(expiration time.Duration, options client.Options) *Client {
	reportsCache := cache.NewCache(expiration)
	c1 := fetcher.NewCacher(
		fetcher.NewLimiter(
			&sciensanoClient.Client{
				Caller: &client.InstrumentedClient{
					Options:     options,
					Application: "sciensano",
				},
			}, 3,
		),
	)
	go c1.AutoRefresh(context.Background(), time.Hour)

	c2 := fetcher.NewCacher(&vaccinesClient.Client{
		Caller: &client.InstrumentedClient{
			Options:     options,
			Application: "vaccines",
		},
	})
	go c2.AutoRefresh(context.Background(), time.Hour)

	return &Client{
		ReportCache:      reportsCache,
		Cases:            cases.Reporter{ReportCache: reportsCache, APIClient: c1},
		Hospitalisations: hospitalisations.Reporter{ReportCache: reportsCache, APIClient: c1},
		Mortality:        mortality.Reporter{ReportCache: reportsCache, APIClient: c1},
		TestResults:      testresults.Reporter{ReportCache: reportsCache, APIClient: c1},
		Vaccinations:     vaccinations.Reporter{ReportCache: reportsCache, APIClient: c1},
		Vaccines:         vaccines.Reporter{ReportCache: reportsCache, APIClient: c2},
	}
}
