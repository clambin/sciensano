package apiclient

import (
	"context"
	"sync"
	"time"
)

// Cache implements a cache for the Sciensano API.  It's meant to be API-compatible with APIClient, so clients can
// replace their APIClient with a Cache instead:
//
// 		client := &Client{
//			APIClient: &apiclient.Client{HTTPClient: &http.Client{}},
//		}
//
// becomes:
//
//		client := &Client{
//			APIClient: &apiclient.Cache{
//				APIClient: &apiclient.Client{HTTPClient: &http.Client{}},
//				Retention: 15 * time.Minute,
//			},
//		}
type Cache struct {
	APIClient
	Retention time.Duration
	lock      sync.Mutex
	entries   map[string]*cacheEntry
}

type cacheEntry struct {
	entries interface{}
	expiry  time.Time
	once    *sync.Once
}

func (cache *Cache) getCacheEntry(name string) (entry *cacheEntry) {
	if cache.entries == nil {
		cache.entries = make(map[string]*cacheEntry)
	}

	var ok bool
	if entry, ok = cache.entries[name]; ok == false {
		entry = &cacheEntry{}
		cache.entries[name] = entry
	}
	return
}

// GetTestResults retrieves all COVID-19 test results.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetTestResults(ctx context.Context) (results []*APITestResultsResponse, err error) {
	cache.lock.Lock()
	entry := cache.getCacheEntry("tests")
	if entry.once == nil || time.Now().After(entry.expiry) {
		entry.once = &sync.Once{}
		entry.expiry = time.Now().Add(cache.Retention)
		metricCacheMiss.WithLabelValues("tests").Add(1.0)
	} else {
		metricCacheHit.WithLabelValues("tests").Add(1.0)
	}
	cache.lock.Unlock()

	entry.once.Do(func() {
		if entry.entries, err = cache.APIClient.GetTestResults(ctx); err == nil {
			cache.entries["tests"] = entry
		}
	})

	return entry.entries.([]*APITestResultsResponse), err
}

// GetVaccinations retrieves all COVID-19 vaccinations.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetVaccinations(ctx context.Context) (results []*APIVaccinationsResponse, err error) {
	cache.lock.Lock()
	entry := cache.getCacheEntry("vaccinations")
	if entry.once == nil || time.Now().After(entry.expiry) {
		entry.once = &sync.Once{}
		entry.expiry = time.Now().Add(cache.Retention)
		metricCacheMiss.WithLabelValues("vaccinations").Add(1.0)
	} else {
		metricCacheHit.WithLabelValues("vaccinations").Add(1.0)
	}
	cache.lock.Unlock()

	entry.once.Do(func() {
		if entry.entries, err = cache.APIClient.GetVaccinations(ctx); err == nil {
			cache.entries["vaccinations"] = entry
		}
	})

	return entry.entries.([]*APIVaccinationsResponse), err
}
