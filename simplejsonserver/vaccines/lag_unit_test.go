package vaccines

import (
	"fmt"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVaccineDelay(t *testing.T) {
	var vaccineData = []struct {
		timestamp time.Time
		partial   float64
		full      float64
	}{
		{timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), partial: 10, full: 0},
		{timestamp: time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC), partial: 15, full: 1},
		{timestamp: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC), partial: 15, full: 4},
		{timestamp: time.Date(2021, 2, 15, 0, 0, 0, 0, time.UTC), partial: 25, full: 5},
		{timestamp: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC), partial: 35, full: 10},
		{timestamp: time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC), partial: 35, full: 15},
	}

	vaccinations := datasets.NewFromAPIResponse(nil)
	for _, entry := range vaccineData {
		vaccinations.Add(entry.timestamp, "partial", entry.partial)
		vaccinations.Add(entry.timestamp, "full", entry.full)
	}

	batches := datasets.NewFromAPIResponse(nil)
	batches.Add(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC), "total", 20)
	batches.Add(time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC), "total", 40)
	batches.Add(time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC), "total", 50)

	expected := []struct {
		Timestamp time.Time
		Value     float64
	}{{
		Timestamp: time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		Value:     45,
	}, {
		Timestamp: time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		Value:     28,
	}, {
		Timestamp: time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		Value:     42,
	}}

	timestamps, delays := calculateVaccineDelay(vaccinations, batches)

	require.Len(t, timestamps, len(expected), timestamps)
	require.Len(t, delays, len(expected))
	for index, entry := range expected {
		assert.Equal(t, entry.Timestamp, timestamps[index], fmt.Sprintf("index: %d", index))
		assert.Equal(t, entry.Value, delays[index], fmt.Sprintf("index: %d", index))
	}
}
