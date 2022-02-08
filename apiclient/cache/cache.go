package cache

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Holder interface contains all functions required to mock a Cache
//go:generate mockery --name Holder
type Holder interface {
	Get(name string) (entries []apiclient.APIResponse, found bool)
	Run(ctx context.Context, interval time.Duration)
	Stats() (stats map[string]int)
}

// Cache holds a list of measurements
type Cache struct {
	lock     sync.RWMutex
	entries  map[string][]apiclient.APIResponse
	Fetchers []Fetcher
}

var _ Holder = &Cache{}

// Get retrieves a cached item
func (c *Cache) Get(name string) (entries []apiclient.APIResponse, found bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	entries, found = c.entries[name]
	return
}

// Run starts the cache and periodically pulls in new data
func (c *Cache) Run(ctx context.Context, interval time.Duration) {
	c.lock.Lock()
	c.entries = make(map[string][]apiclient.APIResponse)
	c.lock.Unlock()

	ch := make(chan FetcherResponse)

	c.startFetchers(ctx, interval, ch)

	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case response := <-ch:
			c.lock.Lock()
			c.entries[response.Name] = response.Response
			c.lock.Unlock()
			log.WithField("name", response.Name).Info("API response cached")
		}
	}
}

func (c *Cache) startFetchers(ctx context.Context, interval time.Duration, ch chan<- FetcherResponse) {
	for _, fetcher := range c.Fetchers {
		go func(f Fetcher) {
			f.Update(ctx, ch)

			ticker := time.NewTicker(interval)
			for running := true; running; {
				select {
				case <-ctx.Done():
					running = false
				case <-ticker.C:
					f.Update(ctx, ch)
				}
			}
			ticker.Stop()
		}(fetcher)
	}
}

// Stats returns statistics on the cache
func (c *Cache) Stats() (stats map[string]int) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	stats = make(map[string]int)
	for key, value := range c.entries {
		stats[key] = len(value)
	}
	return
}
