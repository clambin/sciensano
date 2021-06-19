package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestForecastTests(t *testing.T) {
	predicted, err := predictTests(0)
	assert.Error(t, err)

	History := 28
	predicted, err = predictTests(History)
	assert.NoError(t, err)
	assert.Len(t, predicted, History-predictor.BatchSize+predictor.ForecastSampleCount)

	History = 365
	predicted, err = predictTests(History)
	assert.NoError(t, err)
	assert.Len(t, predicted, predictor.BatchSize*(predictor.HistoryBatches-1)+predictor.ForecastSampleCount)
}

func BenchmarkForecastTests(b *testing.B) {
	const History = 365
	predicted, err := predictTests(History)
	assert.NoError(b, err)
	assert.Len(b, predicted, predictor.BatchSize*(predictor.HistoryBatches-1)+predictor.ForecastSampleCount)
}

func predictTests(history int) ([]sciensano.Test, error) {
	tests := make([]sciensano.Test, 0)
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < history; i++ {
		tests = append(tests, sciensano.Test{
			Timestamp: timestamp,
			Total:     i,
			Positive:  i / 2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}

	return predictor.ForecastTests(tests)
}
