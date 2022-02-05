package datasets_test

import (
	"fmt"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestTimestamps_Add(t *testing.T) {
	ts := datasets.MakeTimestamps()

	for i := 0; i < 100; i++ {
		ts.Add(time.Date(2020, time.January, 1+rand.Intn(31), rand.Intn(24), 0, 0, 0, time.UTC))
	}

	assert.Len(t, ts.List(), ts.Count())

	previous := time.Time{}
	for _, timestamp := range ts.List() {
		assert.True(t, timestamp.After(previous) || timestamp.Equal(previous), fmt.Sprintf("%v, %v", previous, timestamp))
		previous = timestamp
	}
}

func TestTimestamps_GetIndex(t *testing.T) {
	ts := datasets.MakeTimestamps()

	timestamps := make(map[time.Time]int)

	const iterations = 100

	for i := 0; i < iterations; i++ {
		timestamp := time.Date(2020, time.January, 1+rand.Intn(31), rand.Intn(24), 0, 0, 0, time.UTC)

		if index, added := ts.Add(timestamp); added {
			timestamps[timestamp] = index
		}
	}

	assert.Len(t, timestamps, ts.Count())

	for timestamp, index := range timestamps {
		assert.Equal(t, index, ts.GetIndex(timestamp), timestamp)
	}

}
