package vaccines

import (
	"fmt"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVaccineDelay(t *testing.T) {
	vaccinations := &sciensano.Vaccinations{
		Timestamps: []time.Time{
			time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 01, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 02, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		},
		Groups: []sciensano.GroupedVaccinationsEntry{{
			Values: []*sciensano.VaccinationsEntry{
				{Partial: 10, Full: 0},
				{Partial: 15, Full: 1},
				{Partial: 15, Full: 4},
				{Partial: 25, Full: 5},
				{Partial: 35, Full: 10},
				{Partial: 35, Full: 15},
			},
		}},
	}

	batches := []*vaccines.Batch{{
		Date:   vaccines.Time{Time: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 20,
	}, {
		Date:   vaccines.Time{Time: time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 40,
	}, {
		Date:   vaccines.Time{Time: time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 50,
	}}

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
