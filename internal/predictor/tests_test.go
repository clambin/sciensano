package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestForecastTests(t *testing.T) {
	tests := make([]sciensano.Test, 0)

	predicted, score, err := predictor.ForecastTests(tests)
	assert.Error(t, err)

	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 365; i++ {
		tests = append(tests, sciensano.Test{
			Timestamp: timestamp,
			Total:     i,
			Positive:  i / 2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}

	predicted, score, err = predictor.ForecastTests(tests)
	if assert.NoError(t, err) {
		assert.Greater(t, score, 0.99)
		if assert.Len(t, predicted, predictor.ForecastBatches*predictor.BatchSize) {
			start := 365.0
			for i := 0; i < predictor.ForecastBatches*predictor.BatchSize; i++ {
				assert.LessOrEqual(t, 100*math.Abs(start-float64(predicted[i].Total))/start, 15.0, i)
				assert.LessOrEqual(t, 100*math.Abs(start/2-float64(predicted[i].Positive))/start/2, 15.0, i)
				start++
			}
		}
	}
}

func BenchmarkForecastTests(b *testing.B) {
	tests := make([]sciensano.Test, 0)

	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 365; i++ {
		tests = append(tests, sciensano.Test{
			Timestamp: timestamp,
			Total:     i,
			Positive:  i / 2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}

	predicted, score, err := predictor.ForecastTests(tests)
	assert.NoError(b, err)
	assert.Greater(b, score, 0.98)
	assert.Len(b, predicted, predictor.ForecastBatches*predictor.BatchSize)

}
