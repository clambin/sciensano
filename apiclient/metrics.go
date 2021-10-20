package apiclient

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics
var (
	metricRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_api_requests_total",
		Help: "Number of Sciensano API calls made",
	}, []string{"endpoint"})
	metricRequestErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_api_request_errors_total",
		Help: "Number of failed Sciensano API calls",
	}, []string{"endpoint"})
	metricRequestLatency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "sciensano_api_latency",
		Help: "Latency of Sciensano API calls",
	}, []string{"endpoint"})
	metricCacheHit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_cache_hit_total",
		Help: "Number of Sciensano API calls served from cache",
	}, []string{"endpoint"})
	metricCacheMiss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_cache_miss_total",
		Help: "Number of Sciensano API calls not served from cache",
	}, []string{"endpoint"})
)
