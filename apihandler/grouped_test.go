package apihandler

import (
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFilledVaccinations(t *testing.T) {
	vaccinations := []sciensano.Vaccination{
		{
			Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Partial:   10,
			Full:      5,
			Booster:   1,
		},
		{
			Timestamp: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			Partial:   20,
			Full:      10,
			Booster:   2,
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

	filled := getFilledVaccinations(timestamps, vaccinations, groupPartial)

	require.Len(t, filled, len(timestamps))
	assert.Equal(t, 10.0, filled[0])
	assert.Equal(t, 20.0, filled[len(timestamps)-1])

	filled = getFilledVaccinations(timestamps, vaccinations, groupFull)

	require.Len(t, filled, len(timestamps))
	assert.Equal(t, 5.0, filled[0])
	assert.Equal(t, 10.0, filled[len(timestamps)-1])

	filled = getFilledVaccinations(timestamps, vaccinations, groupBooster)

	require.Len(t, filled, len(timestamps))
	assert.Equal(t, 1.0, filled[0])
	assert.Equal(t, 2.0, filled[len(timestamps)-1])
}
