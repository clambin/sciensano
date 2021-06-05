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
	if assert.NoError(t, err) {
		assert.Len(t, predicted, 21)
		start := 365
		for i := 0; i < 7; i++ {
			assert.Less(t, math.Abs(float64(start-predicted[i].Total)), 10.0, i)
			assert.Less(t, math.Abs(float64(start/2-predicted[i].Positive)), 10.0, i)
			start++
		}
	}
}
