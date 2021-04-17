package apihandler

import (
	"fmt"
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVaccineDelay(t *testing.T) {
	vaccinations := []sciensano.Vaccination{{
		Timestamp:  time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
		FirstDose:  10,
		SecondDose: 0,
	}, {
		Timestamp:  time.Date(2021, 01, 15, 0, 0, 0, 0, time.UTC),
		FirstDose:  15,
		SecondDose: 1,
	}, {
		Timestamp:  time.Date(2021, 02, 1, 0, 0, 0, 0, time.UTC),
		FirstDose:  15,
		SecondDose: 4,
	}, {
		Timestamp:  time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		FirstDose:  25,
		SecondDose: 5,
	}, {
		Timestamp:  time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		FirstDose:  35,
		SecondDose: 10,
	}, {
		Timestamp:  time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		FirstDose:  35,
		SecondDose: 15,
	}}

	batches := []vaccines.Batch{{
		Date:   vaccines.Time(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 20,
	}, {
		Date:   vaccines.Time(time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 40,
	}, {
		Date:   vaccines.Time(time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC)),
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

	if assert.Equal(t, len(expected), len(timestamps)) && assert.Equal(t, len(expected), len(delays)) {
		for index, entry := range expected {
			assert.Equal(t, entry.Timestamp, timestamps[index], fmt.Sprintf("index: %d", index))
			assert.Equal(t, entry.Value, delays[index], fmt.Sprintf("index: %d", index))
		}

	}
}
