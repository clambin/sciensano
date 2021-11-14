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
	vaccinations := &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 01, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 02, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "partial", Values: []float64{10, 15, 15, 25, 35, 35}},
			{Name: "full", Values: []float64{0, 1, 4, 5, 10, 15}}},
	}

	batches := &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{{
			Name:   "",
			Values: []float64{20, 40, 50},
		}},
	}

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
