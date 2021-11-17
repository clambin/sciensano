package reporter_test

import (
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := reporter.NewCache(15 * time.Minute)

	for i := 0; i < 10; i++ {
		e := c.Load("foo")
		e.Once.Do(func() {
			if e.Data == nil {
				e.Data = &datasets.Dataset{
					Timestamps: make([]time.Time, 0),
					Groups:     nil,
				}
			}
			e.Data.Timestamps = append(e.Data.Timestamps, time.Now())
		})
		c.Save("foo", e)
	}

	e := c.Load("foo")
	assert.Len(t, e.Data.Timestamps, 1)

	stats := c.Stats()
	assert.Len(t, stats, 1)
}

func TestCache_Stats(t *testing.T) {
	c := reporter.NewCache(100 * time.Millisecond)
	e := c.Load("foo")
	e.Once.Do(func() {
		if e.Data == nil {
			e.Data = &datasets.Dataset{
				Timestamps: make([]time.Time, 0),
				Groups:     nil,
			}
		}
		e.Data.Timestamps = append(e.Data.Timestamps, time.Now())
	})
	c.Save("foo", e)
	stats := c.Stats()
	count, ok := stats["foo"]
	assert.True(t, ok)
	assert.Equal(t, 1, count)

	assert.Eventually(t, func() bool {
		count, ok = c.Stats()["foo"]
		return ok && count == 0
	}, 500*time.Millisecond, 100*time.Millisecond)
}
