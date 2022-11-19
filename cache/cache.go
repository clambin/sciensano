package cache

import (
	"context"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/clambin/sciensano/pkg/limiter"
	"net/http"
	"sync"
	"time"
)

type SciensanoCache struct {
	Cases            cacher[sciensano.Cases]
	Hospitalisations cacher[sciensano.Hospitalisations]
	Mortalities      cacher[sciensano.Mortalities]
	TestResults      cacher[sciensano.TestResults]
	Vaccinations     cacher[sciensano.Vaccinations]
}

func NewSciensanoCache(target string) *SciensanoCache {
	if target == "" {
		target = sciensano.BaseURL
	}

	// TODO: client should be a Limiter too
	client := limiter.NewLimiter(&httpclient.InstrumentedClient{
		BaseClient: httpclient.BaseClient{HTTPClient: http.DefaultClient},
		// TODO: pass prometheus Registerer to metrics. Don't use promauto
		Options:     httpclient.Options{PrometheusMetrics: httpclient.NewMetrics("sciensano", "")},
		Application: "sciendano",
	}, 3)

	return &SciensanoCache{
		Cases: cacher[sciensano.Cases]{
			Fetcher: &fetcher[sciensano.Cases]{
				client: client,
				target: target + sciensano.Routes["cases"],
			},
		},
		Hospitalisations: cacher[sciensano.Hospitalisations]{
			Fetcher: &fetcher[sciensano.Hospitalisations]{
				client: client,
				target: target + sciensano.Routes["hospitalisations"],
			},
		},
		Mortalities: cacher[sciensano.Mortalities]{
			Fetcher: &fetcher[sciensano.Mortalities]{
				client: client,
				target: target + sciensano.Routes["mortalities"],
			},
		},
		TestResults: cacher[sciensano.TestResults]{
			Fetcher: &fetcher[sciensano.TestResults]{
				client: client,
				target: target + sciensano.Routes["testResults"],
			},
		},
		Vaccinations: cacher[sciensano.Vaccinations]{
			Fetcher: &fetcher[sciensano.Vaccinations]{
				client: client,
				target: target + sciensano.Routes["vaccinations"],
			},
		},
	}
}

func (c *SciensanoCache) AutoRefresh(ctx context.Context, interval time.Duration) {
	var wg sync.WaitGroup
	wg.Add(5)
	go func() { defer wg.Done(); c.Cases.AutoRefresh(ctx, interval) }()
	go func() { defer wg.Done(); c.Hospitalisations.AutoRefresh(ctx, interval) }()
	go func() { defer wg.Done(); c.Mortalities.AutoRefresh(ctx, interval) }()
	go func() { defer wg.Done(); c.TestResults.AutoRefresh(ctx, interval) }()
	go func() { defer wg.Done(); c.Vaccinations.AutoRefresh(ctx, interval) }()
	wg.Wait()
}
