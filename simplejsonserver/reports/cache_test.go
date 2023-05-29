package reports_test

import (
	"context"
	"errors"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/simplejsonserver/reports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := reports.NewCache(15 * time.Minute)

	var updates int
	for i := 0; i < 10; i++ {
		e := c.Load("foo")
		e.Once.Do(func() {
			if e.Data == nil {
				e.Data, _ = createBigDataSet()
			}
			c.Save("foo", e)
			updates++
		})
	}

	assert.Equal(t, 1, updates)

	e := c.Load("foo")
	l := len(e.Data.GetTimestamps())
	assert.Equal(t, 500, l)

	stats := c.Stats()
	assert.Len(t, stats, 1)
}

func TestCache_Stats(t *testing.T) {
	c := reports.NewCache(100 * time.Millisecond)
	e := c.Load("foo")
	e.Once.Do(func() {
		if e.Data == nil {
			e.Data, _ = createBigDataSet()
		}
		c.Save("foo", e)
	})

	stats := c.Stats()
	assert.Len(t, stats, 1)

	count, ok := stats["foo"]
	assert.True(t, ok)
	assert.Equal(t, 500, count)

	assert.Eventually(t, func() bool {
		count, ok = c.Stats()["foo"]
		return ok && count == 0
	}, 500*time.Millisecond, 100*time.Millisecond)
}

func TestCache_MaybeGenerate(t *testing.T) {
	called := 0
	c := reports.NewCache(time.Hour)
	_, err := c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
		d, _ := createBigDataSet()
		called++
		return d, nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, called)

	_, err = c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
		d, _ := createBigDataSet()
		called++
		return d, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, called)
}

func TestCache_MaybeGenerate_Stress(t *testing.T) {
	c := reports.NewCache(200 * time.Millisecond)

	const maxParallel = 1e2
	s := semaphore.NewWeighted(maxParallel)
	ctx := context.Background()

	for i := 0; i < 1e4; i++ {
		_ = s.Acquire(ctx, 1)
		go func(i int) {
			report, err := c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
				if rand.Intn(10) < 1 {
					return nil, errors.New("fail")
				}
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				d := tabulator.New("value")
				d.Add(time.Now(), "value", float64(i))
				return d, nil
			})
			assert.False(t, err == nil && report == nil)
			s.Release(1)
		}(i)
	}

	_ = s.Acquire(ctx, maxParallel)
}

func BenchmarkCache_MaybeGenerate(b *testing.B) {
	c := reports.NewCache(time.Second)
	_, err := c.MaybeGenerate("foo", createBigDataSet)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = c.MaybeGenerate("foo", createBigDataSet)
		if err != nil {
			b.Fatal()
		}
	}
}

func createBigDataSet() (d *tabulator.Tabulator, err error) {
	d = tabulator.New("0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
	for r := 0; r < 500; r++ {
		timestamp := time.Now()
		for c := 0; c < 10; c++ {
			d.Add(timestamp, strconv.Itoa(c), float64(r))
		}
	}
	return
}
