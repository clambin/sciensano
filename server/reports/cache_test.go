package reports_test

import (
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/server/reports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestReportCache_MaybeGenerate(t *testing.T) {
	c := reports.ReportCache{Cache: reports.NewLocalCache(time.Minute)}
	var wg sync.WaitGroup
	var tracker concurrencyTracker
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
				tracker.Inc()
				defer tracker.Dec()

				time.Sleep(time.Duration(75+rand.Intn(50)) * time.Millisecond)
				return createSimpleDataSet(), nil
			})
			require.NoError(t, err)
		}()
	}
	wg.Wait()
	assert.Equal(t, 1, tracker.maxInFlight)
}

type concurrencyTracker struct {
	inFlight    int
	maxInFlight int
	lock        sync.Mutex
}

func (c *concurrencyTracker) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inFlight++
	if c.inFlight > c.maxInFlight {
		c.maxInFlight = c.inFlight
	}
}

func (c *concurrencyTracker) Dec() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inFlight--
	if c.inFlight < 0 {
		panic("inFlight cannot be negative")
	}
}

func createBigDataSet() (*tabulator.Tabulator, error) {
	d := tabulator.New("0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
	timestamp := time.Now()
	for r := 0; r < 500; r++ {
		for c := 0; c < 10; c++ {
			d.Add(timestamp, strconv.Itoa(c), float64(r))
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return d, nil
}

func createSimpleDataSet() *tabulator.Tabulator {
	d := tabulator.New("A", "B", "C")
	d.Add(time.Date(2023, time.July, 31, 0, 0, 0, 0, time.UTC), "C", 3)
	d.Add(time.Date(2023, time.July, 30, 0, 0, 0, 0, time.UTC), "C", 2)
	d.Add(time.Date(2023, time.July, 29, 0, 0, 0, 0, time.UTC), "C", 1)
	return d
}
