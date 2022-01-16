package metrics

import (
	"github.com/clambin/go-metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ErrorMetric measures API call failures
	ErrorMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_api_errors_total",
		Help: "Number of failed Sciensano API calls",
	}, []string{"endpoint"})

	// LatencyMetric measures API call duration
	LatencyMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "sciensano_api_latency",
		Help: "Latency of Sciensano API calls",
	}, []string{"endpoint"})

	// Metrics contains all APIClientMetrics
	Metrics = metrics.APIClientMetrics{
		Latency: LatencyMetric,
		Errors:  ErrorMetric,
	}
)
