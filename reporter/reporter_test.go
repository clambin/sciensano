package reporter_test

import (
	"github.com/clambin/go-metrics/tools"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/reporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReporterMetrics(t *testing.T) {
	m := httpclient.NewMetrics("foo", "")
	r := reporter.NewWithOptions(time.Hour, httpclient.Options{PrometheusMetrics: m})

	assert.Eventually(t, func() bool {
		_, err := r.TestResults.Get()
		return err == nil
	}, time.Minute, time.Second)

	len1 := len(r.Cases.APIClient.DataTypes())
	len2 := len(r.Vaccines.APIClient.DataTypes())
	l := len1 + len2

	ch := make(chan prometheus.Metric)
	go m.Latency.Collect(ch)

	for i := 0; i < l; i++ {
		metric := <-ch
		assert.Equal(t, "foo_api_latency", tools.MetricName(metric))
		assert.Contains(t, []string{"sciensano", "vaccines"}, tools.MetricLabel(metric, "application"))
		assert.Contains(t, []string{"GET", "HEAD"}, tools.MetricLabel(metric, "method"))
	}

	ch = make(chan prometheus.Metric)
	go m.Errors.Collect(ch)

	for i := 0; i < l; i++ {
		metric := <-ch
		assert.Equal(t, "foo_api_errors_total", tools.MetricName(metric))
		assert.Contains(t, []string{"sciensano", "vaccines"}, tools.MetricLabel(metric, "application"))
		assert.Contains(t, []string{"GET", "HEAD"}, tools.MetricLabel(metric, "method"))
	}
}
