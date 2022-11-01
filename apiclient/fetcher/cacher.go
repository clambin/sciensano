package fetcher

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"sync"
	"time"
)

type Cacher struct {
	Fetcher
	cache map[int]*cache
}

var _ Fetcher = &Cacher{}

func NewCacher(f Fetcher) *Cacher {
	caches := make(map[int]*cache)
	for dataType := range f.DataTypes() {
		caches[dataType] = &cache{
			fetcher:  f,
			dataType: dataType,
			expiry:   time.Hour,
		}
	}
	return &Cacher{
		Fetcher: f,
		cache:   caches,
	}
}

func (c *Cacher) Fetch(ctx context.Context, dataType int) ([]apiclient.APIResponse, error) {
	dataTypeCache, ok := c.cache[dataType]
	if !ok {
		return nil, fmt.Errorf("invalid data type: %d", dataType)
	}
	if err := dataTypeCache.refresh(ctx); err != nil {
		return nil, err
	}
	return dataTypeCache.Get(), nil
}

func (c *Cacher) AutoUpdate(ctx context.Context, interval time.Duration) {
	var wg sync.WaitGroup
	wg.Add(len(c.cache))
	for _, dataTypeCache := range c.cache {
		go func(cache *cache) { cache.Run(ctx, interval); wg.Done() }(dataTypeCache)
	}
	<-ctx.Done()
	wg.Wait()
}
