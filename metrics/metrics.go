package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
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
)

type TimerMetric struct {
	timer *prometheus.Timer
	name  string
}

func NewTimerMetric(name string) *TimerMetric {
	return &TimerMetric{
		name:  name,
		timer: prometheus.NewTimer(metricRequestLatency.WithLabelValues(name)),
	}
}

func (tm TimerMetric) Report(pass bool) {
	duration := tm.timer.ObserveDuration()
	log.WithField("duration", duration).Debugf("called %s API", tm.name)
	metricRequestsTotal.WithLabelValues(tm.name).Add(1.0)
	if pass == false {
		metricRequestErrorsTotal.WithLabelValues(tm.name).Add(1.0)
	}
}
