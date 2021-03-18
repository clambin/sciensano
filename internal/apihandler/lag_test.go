package apihandler

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVaccinationLag(t *testing.T) {
	vaccinations := []sciensano.Vaccination{
		{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), FirstDose: 0, SecondDose: 0},
		{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), FirstDose: 1, SecondDose: 0},
		{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), FirstDose: 2, SecondDose: 1},
		{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), FirstDose: 3, SecondDose: 2},
		{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), FirstDose: 4, SecondDose: 3},
		{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), FirstDose: 5, SecondDose: 4},
		{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), FirstDose: 6, SecondDose: 5},
	}
	lag := buildLag(vaccinations)

	if assert.Len(t, lag, 5) {
		assert.Equal(t, 1.0, lag[0].Lag)
		assert.Equal(t, 1.0, lag[1].Lag)
		assert.Equal(t, 1.0, lag[2].Lag)
		assert.Equal(t, 1.0, lag[3].Lag)
	}

	vaccinations = []sciensano.Vaccination{
		{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), FirstDose: 1, SecondDose: 1}, // 0
		{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), FirstDose: 1, SecondDose: 1}, // -
		{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), FirstDose: 2, SecondDose: 1}, // 1
		{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), FirstDose: 3, SecondDose: 2}, // 1
		{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), FirstDose: 4, SecondDose: 3}, // 1
		{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), FirstDose: 5, SecondDose: 4}, // 1
		{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), FirstDose: 6, SecondDose: 4}, // 2
	}

	lag = buildLag(vaccinations)

	if assert.Len(t, lag, 6) {
		assert.Equal(t, 0.0, lag[0].Lag)
		assert.Equal(t, 1.0, lag[1].Lag)
		assert.Equal(t, 1.0, lag[2].Lag)
		assert.Equal(t, 1.0, lag[3].Lag)
		assert.Equal(t, 1.0, lag[4].Lag)
		assert.Equal(t, 2.0, lag[5].Lag)
	}

}
