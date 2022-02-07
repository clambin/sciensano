package cache_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

type testFetcher struct {
	count int
}

func (t *testFetcher) call(_ context.Context) (response []apiclient.APIResponse, err error) {
	if t.count < 2 {
		t.count++
		return nil, errors.New("not yet ...")
	}
	return response, nil
}

func (t *testFetcher) Update(ctx context.Context, ch chan<- cache.FetcherResponse) {
	cache.Fetch(ctx, ch, "test", t.call)
}

var _ cache.Fetcher = &testFetcher{}

func TestFetch(t *testing.T) {
	f := &testFetcher{}
	c := cache.Cache{Fetchers: []cache.Fetcher{f}}

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.Run(ctx, time.Hour)
		wg.Done()
	}()

	require.Eventually(t, func() bool {
		_, ok := c.Get("test")
		return ok
	}, time.Minute, 10*time.Millisecond)

	cancel()
	wg.Wait()

}
