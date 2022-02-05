package vaccinations

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnitVaccinationLag(t *testing.T) {
	vaccinations := datasets.New()

	for _, inputData := range []struct {
		timestamp     time.Time
		partial, full float64
	}{
		{timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), partial: 0, full: 0},
		{timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), partial: 1, full: 0},
		{timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), partial: 2, full: 1},
		{timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), partial: 3, full: 2},
		{timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), partial: 4, full: 3},
		{timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), partial: 5, full: 4},
		{timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), partial: 6, full: 5},
	} {
		vaccinations.Add(inputData.timestamp, "partial", inputData.partial)
		vaccinations.Add(inputData.timestamp, "full", inputData.full)
	}

	_, lag := buildLag(vaccinations)
	assert.Equal(t, query.NumberColumn{1.0, 1.0, 1.0, 1.0, 1.0}, lag)

	vaccinations = datasets.New()
	for _, inputData := range []struct {
		timestamp     time.Time
		partial, full float64
	}{
		{timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), partial: 1, full: 1},
		{timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), partial: 1, full: 1},
		{timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), partial: 2, full: 1},
		{timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), partial: 3, full: 2},
		{timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), partial: 4, full: 3},
		{timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), partial: 5, full: 4},
		{timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), partial: 6, full: 5},
	} {
		vaccinations.Add(inputData.timestamp, "partial", inputData.partial)
		vaccinations.Add(inputData.timestamp, "full", inputData.full)
	}

	_, lag = buildLag(vaccinations)
	assert.Equal(t, query.NumberColumn{0.0, 1.0, 1.0, 1.0, 1.0}, lag)
}
