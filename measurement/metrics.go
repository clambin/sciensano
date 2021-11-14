package measurement

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics
var (
	metricCacheTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_cache_hit_total",
		Help: "Total API calls attempted from cache",
	}, []string{"endpoint"})
	metricCacheMiss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_cache_miss_total",
		Help: "Number of API calls not served from cache",
	}, []string{"endpoint"})
)
