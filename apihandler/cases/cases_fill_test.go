package cases

import (
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFillCases(t *testing.T) {
	input := map[string][]sciensano.CaseCount{
		"A": {
			{Timestamp: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC), Count: 10},
			{Timestamp: time.Date(2021, 10, 25, 0, 0, 0, 0, time.UTC), Count: 20},
			{Timestamp: time.Date(2021, 10, 26, 0, 0, 0, 0, time.UTC), Count: 30},
		},
		"B": {
			{Timestamp: time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC), Count: 15},
		},
	}

	output := fillCases(input)

	expected := map[time.Time]map[string]int{
		time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC): {
			"A": 10,
			"B": 0,
		},
		time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC): {
			"A": 0,
			"B": 15,
		},
		time.Date(2021, 10, 25, 0, 0, 0, 0, time.UTC): {
			"A": 20,
			"B": 0,
		},
		time.Date(2021, 10, 26, 0, 0, 0, 0, time.UTC): {
			"A": 30,
			"B": 0,
		},
	}

	assert.Equal(t, expected, output)
}
