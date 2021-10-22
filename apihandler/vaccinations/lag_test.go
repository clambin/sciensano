package vaccinations

import (
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVaccinationLag(t *testing.T) {
	vaccinations := []sciensano.Vaccination{
		{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Partial: 0, Full: 0},
		{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Partial: 1, Full: 0},
		{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Partial: 2, Full: 1},
		{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Partial: 3, Full: 2},
		{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), Partial: 4, Full: 3},
		{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), Partial: 5, Full: 4},
		{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), Partial: 6, Full: 5},
	}
	_, lag := buildLag(vaccinations)

	if assert.Len(t, lag, 5) {
		assert.Equal(t, 1.0, lag[0])
		assert.Equal(t, 1.0, lag[1])
		assert.Equal(t, 1.0, lag[2])
		assert.Equal(t, 1.0, lag[3])
	}

	vaccinations = []sciensano.Vaccination{
		{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Partial: 1, Full: 1}, // 0
		{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Partial: 1, Full: 1}, // -
		{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Partial: 2, Full: 1}, // -
		{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Partial: 3, Full: 2}, // 1
		{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), Partial: 4, Full: 3}, // 1
		{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), Partial: 4, Full: 4}, // 1
		{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), Partial: 6, Full: 5}, // 0
	}

	_, lag = buildLag(vaccinations)

	if assert.Len(t, lag, 5) {
		assert.Equal(t, 0.0, lag[0])
		assert.Equal(t, 1.0, lag[1])
		assert.Equal(t, 1.0, lag[2])
		assert.Equal(t, 1.0, lag[3])
		assert.Equal(t, 0.0, lag[4])
	}

}
