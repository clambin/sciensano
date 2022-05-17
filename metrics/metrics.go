package metrics

import (
	"github.com/clambin/go-metrics/caller"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ErrorMetric measures API call failures
	ErrorMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_api_errors_total",
		Help: "Number of failed Reporter API calls",
	}, []string{"application", "endpoint"})

	// LatencyMetric measures API call duration
	LatencyMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "sciensano_api_latency",
		Help: "Latency of Reporter API calls",
	}, []string{"application", "endpoint"})

	// Metrics contains all APIClientMetrics
	Metrics = caller.ClientMetrics{
		Latency: LatencyMetric,
		Errors:  ErrorMetric,
	}
)
