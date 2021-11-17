package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics
var (
	MetricRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_api_requests_total",
		Help: "Number of Sciensano API calls made",
	}, []string{"endpoint"})
	MetricRequestErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_api_request_errors_total",
		Help: "Number of failed Sciensano API calls",
	}, []string{"endpoint"})
	MetricRequestLatency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "sciensano_api_latency",
		Help: "Latency of Sciensano API calls",
	}, []string{"endpoint"})
)
