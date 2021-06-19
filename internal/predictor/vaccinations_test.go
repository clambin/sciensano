package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestForecastVaccinations(t *testing.T) {
	predicted, err := predictVaccinations(0)
	assert.Error(t, err)

	predicted, err = predictVaccinations(28)
	assert.NoError(t, err)
	assert.Len(t, predicted, 28-predictor.BatchSize+predictor.ForecastSampleCount)
}

func BenchmarkForecastVaccinations(b *testing.B) {
	predicted, err := predictVaccinations(365)
	assert.NoError(b, err)
	assert.Len(b, predicted, 365-predictor.BatchSize+predictor.ForecastSampleCount)
}

func predictVaccinations(history int) ([]sciensano.Vaccination, error) {
	input := make([]sciensano.Vaccination, 0)
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < history; i++ {
		input = append(input, sciensano.Vaccination{
			Timestamp:  timestamp,
			FirstDose:  i,
			SecondDose: i / 2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}

	return predictor.ForecastVaccinations(input)
}
