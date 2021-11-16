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
}

// Cache holds a list of measurements
type Cache struct {
	lock    sync.RWMutex
	entries map[string][]Measurement
	Fetcher Fetcher
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
	if entries, err := c.Fetcher.Update(ctx); err == nil {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.entries = entries
	}
}

// CacheSize returns the number of entries currently in the cache
func (c *Cache) CacheSize() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.entries)
}
