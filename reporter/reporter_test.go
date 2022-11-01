package reporter_test

import (
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/reporter"
	"github.com/prometheus/client_golang/prometheus"
	pcg "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	metrics, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)
	var latencyCount, errorCount int
	for _, metric := range metrics {
		//for range metric.Metric {
		switch *metric.Name {
		case "foo_api_latency":
			assert.Equal(t, pcg.MetricType_SUMMARY, *metric.Type)
			//assert.Equal(t, uint64(1), entry.Summary.GetSampleCount())
			//assert.NotZero(t, entry.Summary.GetSampleSum())
			latencyCount++
		case "foo_api_errors_total":
			assert.Equal(t, pcg.MetricType_COUNTER, *metric.Type)
			//assert.Zero(t, entry.Counter.GetValue())
			errorCount++
		}
		//}
	}
	assert.NotZero(t, latencyCount)
	assert.NotZero(t, errorCount)
}
