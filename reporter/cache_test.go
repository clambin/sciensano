package reporter_test

import (
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3/dataset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := reporter.NewCache(15 * time.Minute)

	for i := 0; i < 10; i++ {
		e := c.Load("foo")
		e.Once.Do(func() {
			if e.Data == nil {
				e.Data = dataset.New()
				e.Data.Add(time.Now(), "A", 1)
			}
			c.Save("foo", e)
		})
	}

	e := c.Load("foo")
	assert.Equal(t, 1, e.Data.Size())

	stats := c.Stats()
	assert.Len(t, stats, 1)
}

func TestCache_Stats(t *testing.T) {
	c := reporter.NewCache(100 * time.Millisecond)
	e := c.Load("foo")
	e.Once.Do(func() {
		if e.Data == nil {
			e.Data = dataset.New()
		}
		e.Data.Add(time.Now(), "A", 1)
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

func TestCache_MaybeGenerate(t *testing.T) {
	called := 0
	c := reporter.NewCache(time.Second)
	result, err := c.MaybeGenerate("foo", func() (*dataset.Dataset, error) {
		d := dataset.New()
		d.Add(time.Now(), "A", 1)
		called++
		return d, nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, called)

	result, err = c.MaybeGenerate("foo", func() (*dataset.Dataset, error) {
		d := dataset.New()
		d.Add(time.Now(), "A", 1)
		called++
		return d, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 1, result.Size())
	require.Equal(t, []string{"A"}, result.GetColumns())
	values, ok := result.GetValues("A")
	require.True(t, ok)
	assert.Len(t, values, 1)
}

func BenchmarkCache_MaybeGenerate(b *testing.B) {
	c := reporter.NewCache(time.Second)
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

func createBigDataSet() (d *dataset.Dataset, err error) {
	d = dataset.New()
	for r := 0; r < 500; r++ {
		for c := 0; c < 10; c++ {
			d.Add(time.Now(), strconv.Itoa(c), 1)
		}
	}
	return
}
