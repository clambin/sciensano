package cache_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

type fetcher struct {
	called int
}

var _ cache.Fetcher = &fetcher{}

type value struct {
	Timestamp time.Time
	Value     float64
}

var _ apiclient.APIResponse = &value{}

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

func (f *fetcher) Fetch(_ context.Context, ch chan<- cache.FetcherResponse) {
	f.called++
	ch <- cache.FetcherResponse{
		Name:     "foo",
		Response: []apiclient.APIResponse{&value{}},
	}
}

func TestCache(t *testing.T) {
	f := &fetcher{}
	ct := &cache.Cache{
		Fetchers: []cache.Fetcher{f},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ct.Run(ctx, 20*time.Millisecond)
		wg.Done()
	}()

	require.Eventually(t, func() bool { return len(ct.Stats()) == 1 }, 500*time.Millisecond, 10*time.Millisecond)

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
