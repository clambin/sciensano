package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestForecastVaccinations(t *testing.T) {
	var score float64
	input := make([]sciensano.Vaccination, 0)

	predicted, score, err := predictor.ForecastVaccinations(input)
	assert.Error(t, err)

	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 365; i++ {
		input = append(input, sciensano.Vaccination{
			Timestamp:  timestamp,
			FirstDose:  i,
			SecondDose: i / 2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}

	predicted, score, err = predictor.ForecastVaccinations(input)
	if assert.NoError(t, err) {
		assert.Greater(t, score, 0.9)
		assert.Len(t, predicted, predictor.ForecastBatches*predictor.BatchSize)
		start := 365
		for i := 0; i < predictor.ForecastBatches*predictor.BatchSize; i++ {
			assert.Less(t, math.Abs(100*float64(start-predicted[i].FirstDose)/float64(start)), 20.0, i)
			assert.Less(t, math.Abs(100*float64(start/2-predicted[i].SecondDose)/float64(start/2)), 20.0, i)
			start++
		}
	}
}

func BenchmarkForecastVaccinations(b *testing.B) {
	input := make([]sciensano.Vaccination, 0)
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 365; i++ {
		input = append(input, sciensano.Vaccination{
			Timestamp:  timestamp,
			FirstDose:  i,
			SecondDose: i / 2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}

	predicted, score, err := predictor.ForecastVaccinations(input)
	if assert.NoError(b, err) {
		assert.Greater(b, score, 0.9)
		assert.Len(b, predicted, predictor.ForecastBatches*predictor.BatchSize)
	}
}
