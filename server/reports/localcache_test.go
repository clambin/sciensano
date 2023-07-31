package reports_test

import (
	"errors"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/server/reports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLocalCache(t *testing.T) {
	c := reports.ReportCache{Cache: reports.NewLocalCache(15 * time.Minute)}
	refTable, _ := createSimpleDataSet()

	var updates int
	for i := 0; i < 10; i++ {
		fromCache, err := c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
			updates++
			return createSimpleDataSet()
		})
		require.NoError(t, err)
		assert.Equal(t, refTable.Size(), fromCache.Size())
		assert.Equal(t, refTable.GetColumns(), fromCache.GetColumns())
	}

	assert.Equal(t, 1, updates)
}

func BenchmarkLocalCache(b *testing.B) {
	c := reports.ReportCache{Cache: reports.NewLocalCache(15 * time.Minute)}
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

func TestCache_Stats(t *testing.T) {
	c := reports.ReportCache{Cache: reports.NewLocalCache(100 * time.Millisecond)}
	_ = c.Stats()
	_, err := c.MaybeGenerate("foo", createSimpleDataSet)
	require.NoError(t, err)

	stats := c.Stats()
	assert.Len(t, stats, 1)

	count, ok := stats["foo"]
	assert.True(t, ok)
	assert.Equal(t, 3, count)

	assert.Eventually(t, func() bool {
		count, ok = c.Stats()["foo"]
		return !ok || count == 0
	}, 500*time.Millisecond, 100*time.Millisecond)
}

func TestLocalCache_Stress(t *testing.T) {
	c := reports.ReportCache{Cache: reports.NewLocalCache(200 * time.Millisecond)}

	const maxParallel = 1e3
	sem := make(chan struct{}, maxParallel)

	var wg sync.WaitGroup

	for i := 0; i < 1e5; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			report, err := c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
				if rand.Intn(10) < 1 {
					return nil, errors.New("fail")
				}
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				d := tabulator.New("value")
				d.Add(time.Now(), "value", float64(i))
				return d, nil
			})
			assert.True(t, report != nil || err != nil)
			<-sem
		}(i)
	}

	wg.Wait()
}
