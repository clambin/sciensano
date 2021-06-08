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
		assert.Len(t, predicted, 28)
		start := 365
		for i := 0; i < 28; i++ {
			assert.Less(t, math.Abs(float64(start-predicted[i].FirstDose)), 5.0, i)
			assert.Less(t, math.Abs(float64(start/2-predicted[i].SecondDose)), 2.5, i)
			start++
		}
	}
}
