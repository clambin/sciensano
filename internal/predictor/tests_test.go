package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestForecastTests(t *testing.T) {
	tests := make([]sciensano.Test, 0)

	predicted, err := predictor.ForecastTests(tests)
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

	predicted, err = predictor.ForecastTests(tests)
	assert.NoError(t, err)
	assert.Len(t, predicted, 365-predictor.BatchSize+predictor.ForecastSamples)
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

	predicted, err := predictor.ForecastTests(tests)
	assert.NoError(b, err)
	assert.Len(b, predicted, predictor.ForecastSamples)
}
