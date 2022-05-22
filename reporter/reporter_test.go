package reporter_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/go-metrics/tools"
	"github.com/clambin/sciensano/reporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReporterMetrics(t *testing.T) {
	m := client.NewMetrics("foo", "")
	r := reporter.NewWithOptions(time.Hour, client.Options{PrometheusMetrics: m})
	go r.APICache.Run(context.Background(), time.Minute)

	assert.Eventually(t, func() bool {
		_, err := r.GetTestResults()
		return err == nil
	}, time.Minute, time.Second)

	ch := make(chan prometheus.Metric)
	go m.Latency.Collect(ch)

	metric := <-ch
	assert.Equal(t, `foo_api_latency`, tools.MetricName(metric))
	//assert.Equal(t, uint64(1), tools.MetricValue(metric).GetSummary().GetSampleCount())
	assert.Equal(t, "GET", tools.MetricLabel(metric, "method"))

	ch = make(chan prometheus.Metric)
	go m.Errors.Collect(ch)

	metric = <-ch
	assert.Equal(t, `foo_api_errors_total`, tools.MetricName(metric))
	//assert.Equal(t, float64(0), tools.MetricValue(metric).GetCounter().GetValue())
	assert.Equal(t, "GET", tools.MetricLabel(metric, "method"))
}
