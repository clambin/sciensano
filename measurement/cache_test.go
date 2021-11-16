package measurement_test

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

type fetcher struct {
	called int
}

type value struct {
	Timestamp time.Time
	Value     float64
}

func (v value) GetTimestamp() time.Time {
	panic("implement me")
}

func (v value) GetGroupFieldValue(_ int) string {
	panic("implement me")
}

func (v value) GetTotalValue() float64 {
	panic("implement me")
}

func (v value) GetAttributeNames() []string {
	panic("implement me")
}

func (v value) GetAttributeValues() []float64 {
	panic("implement me")
}

func (f *fetcher) Update(_ context.Context) (entries map[string][]measurement.Measurement, err error) {
	f.called++
	entries = map[string][]measurement.Measurement{
		"foo": {
			&value{},
		},
	}
	return
}

func TestCache(t *testing.T) {
	f := &fetcher{}
	ct := &measurement.Cache{
		Fetcher: f,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ct.AutoRefresh(ctx, 20*time.Millisecond)
		wg.Done()
	}()

	require.Eventually(t, func() bool { return ct.CacheSize() == 1 }, 500*time.Millisecond, 10*time.Millisecond)

	entries, found := ct.Get("foo")
	require.True(t, found)
	assert.Len(t, entries, 1)

	_, found = ct.Get("bar")
	require.False(t, found)

	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	assert.Greater(t, f.called, 1)

}
