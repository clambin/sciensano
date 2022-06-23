package fetcher

import (
	"context"
	"fmt"
	"github.com/clambin/cache"
	"github.com/clambin/sciensano/apiclient"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Cacher struct {
	GracePeriod time.Duration
	Expiration  time.Duration
	Fetcher
	cache.Cache[int, cacheEntry]
	lock sync.Mutex
}

var _ Fetcher = &Cacher{}

func NewCacher(f Fetcher) *Cacher {
	return &Cacher{
		GracePeriod: 15 * time.Minute,
		Expiration:  time.Hour,
		Fetcher:     f,
		Cache:       *cache.New[int, cacheEntry](0, 0),
	}
}

func (c *Cacher) Fetch(ctx context.Context, dataType int) (results []apiclient.APIResponse, err error) {
	var entry cacheEntry
	entry, err = c.checkCache(ctx, dataType)
	if err != nil || entry.once == nil {
		results = entry.entries
		return
	}

	entry.once.Do(func() {
		log.Debugf("getting new data for %s", c.DataTypes()[dataType])
		if results, err = c.Fetcher.Fetch(ctx, dataType); err == nil {
			c.Cache.Add(dataType, cacheEntry{
				entries: results,
				updated: time.Now(),
				checked: time.Now(),
			})
		}
		log.WithError(err).Debugf("received %d entries for %s", len(results), c.DataTypes()[dataType])
	})

	if err != nil {
		return
	}

	var found bool
	entry, found = c.Cache.Get(dataType)
	if !found {
		err = fmt.Errorf("data for %d not available", dataType)
	}
	return entry.entries, err
}

func (c *Cacher) GetLastUpdates(ctx context.Context, dataType int) (time.Time, error) {
	return c.Fetcher.GetLastUpdates(ctx, dataType)
}

func (c *Cacher) AutoRefresh(ctx context.Context, interval time.Duration) {
	c.refresh(ctx)
	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			c.refresh(ctx)
		}
	}
	ticker.Stop()
}

func (c *Cacher) refresh(ctx context.Context) {
	wg := sync.WaitGroup{}
	for dataType, name := range c.DataTypes() {
		wg.Add(1)
		go func(d int, n string) {
			if _, err := c.Fetch(ctx, d); err != nil {
				log.WithError(err).Errorf("failed to update %s", n)
			}
			wg.Done()
		}(dataType, name)
	}
	wg.Wait()
}

func (c *Cacher) checkCache(ctx context.Context, dataType int) (entry cacheEntry, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var lastModified time.Time
	var found bool

	// do we already have this dataType in cache?
	entry, found = c.Cache.Get(dataType)

	if found {
		// do we need to check when the data was last updated?
		if entry.shouldPoll(c.GracePeriod, c.Expiration) {
			// Is there more recent data?
			if lastModified, err = c.GetLastUpdates(ctx, dataType); err != nil {
				err = fmt.Errorf("failed to get LastModified: %w", err)
				return
			}
		}
	}

	// not in cache, or newer data is available, set up Once to retrieve the data
	if !found || lastModified.After(entry.updated) {
		entry.once = &sync.Once{}
	} else {
		entry.checked = time.Now()
	}
	c.Cache.Add(dataType, entry)

	return
}
