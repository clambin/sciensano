package apihandler

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFilledVaccinations(t *testing.T) {
	vaccinations := []sciensano.Vaccination{
		{
			Timestamp:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			FirstDose:  10,
			SecondDose: 5,
		},
		{
			Timestamp:  time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			FirstDose:  20,
			SecondDose: 10,
		},
	}
	timestamps := []time.Time{
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
	}

	filled := getFilledVaccinations(timestamps, vaccinations, false)

	if assert.Len(t, filled, len(timestamps)) {
		assert.Equal(t, 10.0, filled[0])
		assert.Equal(t, 20.0, filled[len(timestamps)-1])
	}

	filled = getFilledVaccinations(timestamps, vaccinations, true)

	if assert.Len(t, filled, len(timestamps)) {
		assert.Equal(t, 5.0, filled[0])
		assert.Equal(t, 10.0, filled[len(timestamps)-1])
	}

}
