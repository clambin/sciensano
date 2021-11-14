package measurement

import (
	"context"
	"sync"
	"time"
)

// Cache holds a list of measurements
type Cache struct {
	Retention time.Duration
	lock      sync.Mutex
	entries   map[string]*CacheEntry
}

// CacheEntry represents one list of measurements in a Cache
type CacheEntry struct {
	Entries []Measurement
	Expiry  time.Time
	Once    *sync.Once
}

// CacheSize returns the number of entries currently in the cache
func (cache *Cache) CacheSize() int {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	return len(cache.entries)
}

// GetCacheEntry retrieves a cached measurement.  If the entry does not exist, it sets up a Once field so the called
// can retrieve the data exactly once.
func (cache *Cache) GetCacheEntry(name string) (entry *CacheEntry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.entries == nil {
		cache.entries = make(map[string]*CacheEntry)
	}

	var ok bool
	if entry, ok = cache.entries[name]; ok == false {
		entry = &CacheEntry{}
	}

	if entry.Once == nil || time.Now().After(entry.Expiry) {
		entry.Once = &sync.Once{}
		entry.Expiry = time.Now().Add(cache.Retention)
		cache.entries[name] = entry
		metricCacheMiss.WithLabelValues(name).Add(1.0)
	}
	metricCacheTotal.WithLabelValues(name).Add(1.0)

	return
}

// SetCacheEntry saves an entry in the Cache
func (cache *Cache) SetCacheEntry(name string, entry *CacheEntry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.entries[name] = entry
}

// Call is a convenience wrapper for Cache. It attempts to returned cache entry. If the entry does not exist,
// fetcher will be called exactly once and the response is saved in the cache.
func (cache *Cache) Call(ctx context.Context, name string, fetcher func(ctx context.Context) (results []Measurement, err error)) (results []Measurement, err error) {
	entry := cache.GetCacheEntry(name)
	entry.Once.Do(func() {
		entry.Entries, err = fetcher(ctx)
		if err != nil {
			entry.Once = nil
		}
		cache.SetCacheEntry(name, entry)
	})

	return entry.Entries, err
}
