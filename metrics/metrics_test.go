package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	pcg "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTimerMetric_Report(t *testing.T) {
	timer := NewTimerMetric("foo")
	timer.Report(false)

	ch := make(chan prometheus.Metric)
	go metricRequestLatency.Collect(ch)
	val := <-ch

	m := validateMetric(t, val, "foo")
	assert.Equal(t, uint64(1), m.GetSummary().GetSampleCount())

	ch = make(chan prometheus.Metric)
	go metricRequestErrorsTotal.Collect(ch)
	val = <-ch

	m = validateMetric(t, val, "foo")
	assert.Equal(t, float64(1), m.GetCounter().GetValue())

	ch = make(chan prometheus.Metric)
	go metricRequestsTotal.Collect(ch)
	val = <-ch

	m = validateMetric(t, val, "foo")
	assert.Equal(t, float64(1), m.GetCounter().GetValue())
}

func validateMetric(t *testing.T, m prometheus.Metric, label string) (metric *pcg.Metric) {
	metric = &pcg.Metric{}
	err := m.Write(metric)
	require.NoError(t, err)
	labels := metric.GetLabel()
	require.Len(t, labels, 1)
	assert.Equal(t, "endpoint", labels[0].GetName())
	assert.Equal(t, label, labels[0].GetValue())
	return
}
