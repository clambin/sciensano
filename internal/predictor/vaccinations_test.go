package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestForecastVaccinations(t *testing.T) {
	input := make([]sciensano.Vaccination, 0)

	predicted, err := predictor.ForecastVaccinations(input)
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

	predicted, err = predictor.ForecastVaccinations(input)
	assert.NoError(t, err)
	assert.Len(t, predicted, 365-predictor.BatchSize+predictor.ForecastSamples)
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

	predicted, err := predictor.ForecastVaccinations(input)
	assert.NoError(b, err)
	assert.Len(b, predicted, 365-predictor.BatchSize+predictor.ForecastSamples)
}
