package measurement

import (
	"context"
	"sync"
	"time"
)

// Fetcher interface contains all functions an API interface needs to implement to be used by Cache
type Fetcher interface {
	Update(ctx context.Context) (entries map[string][]Measurement, err error)
}

// Holder interface contains all functions required to mock a Cache
//go:generate mockery --name Holder
type Holder interface {
	Get(name string) (entries []Measurement, found bool)
	AutoRefresh(ctx context.Context, interval time.Duration)
	Refresh(ctx context.Context)
	Stats() (stats map[string]int)
}

// Cache holds a list of measurements
type Cache struct {
	lock     sync.RWMutex
	entries  map[string][]Measurement
	Fetchers []Fetcher
}

// Get retrieves a cached item
func (c *Cache) Get(name string) (entries []Measurement, found bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	entries, found = c.entries[name]
	return
}

// AutoRefresh periodically updates the cache
func (c *Cache) AutoRefresh(ctx context.Context, interval time.Duration) {
	c.Refresh(ctx)

	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			c.Refresh(ctx)
		}
	}
}

// Refresh updates the cache
func (c *Cache) Refresh(ctx context.Context) {
	newEntries := make(map[string][]Measurement)
	for _, fetcher := range c.Fetchers {
		if entries, err := fetcher.Update(ctx); err == nil {
			for key, value := range entries {
				newEntries[key] = value
			}
		}
	}
	if len(newEntries) > 0 {
		c.lock.Lock()
		c.entries = newEntries
		c.lock.Unlock()
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
