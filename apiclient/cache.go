package apiclient

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Cache implements a cache for the Sciensano API.  It's meant to be API-compatible with Getter, so clients can
// replace their Getter with a Cache instead:
//
// 		client := &Client{
//			Getter: &apiclient.Client{HTTPClient: &http.Client{}},
//		}
//
// becomes:
//
//		client := &Client{
//			Getter: &apiclient.Cache{
//				Getter: &apiclient.Client{HTTPClient: &http.Client{}},
//				Retention: 15 * time.Minute,
//			},
//		}
type Cache struct {
	Getter
	Retention time.Duration
	lock      sync.Mutex
	entries   map[string]*cacheEntry
}

type cacheEntry struct {
	entries []Measurement
	expiry  time.Time
	once    *sync.Once
}

// AutoRefresh refreshes the cache on a period basis
func (cache *Cache) AutoRefresh(ctx context.Context) {
	cache.refresh(ctx)
	ticker := time.NewTicker(cache.Retention)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			cache.refresh(ctx)
		}
	}
}

func (cache *Cache) refresh(ctx context.Context) {
	before := time.Now()
	log.Debug("refreshing API cache")

	_, err := cache.GetVaccinations(ctx)
	if err != nil {
		log.WithError(err).Warning("failed to update vaccinations cache")
	}
	_, err = cache.GetCases(ctx)
	if err != nil {
		log.WithError(err).Warning("failed to update cases cache")
	}
	_, err = cache.GetHospitalisations(ctx)
	if err != nil {
		log.WithError(err).Warning("failed to update hospitalisations cache")
	}
	_, err = cache.GetMortality(ctx)
	if err != nil {
		log.WithError(err).Warning("failed to update mortality cache")
	}
	_, err = cache.GetTestResults(ctx)
	if err != nil {
		log.WithError(err).Warning("failed to update test results cache")
	}

	log.WithFields(log.Fields{"duration": time.Now().Sub(before), "cache": cache.CacheSize()}).Debugf("refreshed API cache")
}

// CacheSize returns the number of entries currently in the cache
func (cache *Cache) CacheSize() int {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	return len(cache.entries)
}

func (cache *Cache) getCacheEntry(name string) (entry *cacheEntry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.entries == nil {
		cache.entries = make(map[string]*cacheEntry)
	}

	var ok bool
	if entry, ok = cache.entries[name]; ok == false {
		entry = &cacheEntry{}
	}

	if entry.once == nil || time.Now().After(entry.expiry) {
		entry.once = &sync.Once{}
		entry.expiry = time.Now().Add(cache.Retention)
		cache.entries[name] = entry
		metricCacheMiss.WithLabelValues(name).Add(1.0)
	} else {
		metricCacheHit.WithLabelValues(name).Add(1.0)
	}

	return
}

func (cache *Cache) setCacheEntry(name string, entry *cacheEntry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.entries[name] = entry
}

// GetTestResults retrieves all COVID-19 test results.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetTestResults(ctx context.Context) (results []Measurement, err error) {
	entry := cache.getCacheEntry("tests")
	entry.once.Do(func() {
		entry.entries, err = cache.Getter.GetTestResults(ctx)
		if err != nil {
			entry.once = nil
		}
		cache.setCacheEntry("tests", entry)
	})

	return entry.entries, err
}

// GetVaccinations retrieves all COVID-19 vaccinations.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetVaccinations(ctx context.Context) (results []Measurement, err error) {
	entry := cache.getCacheEntry("vaccinations")
	entry.once.Do(func() {
		entry.entries, err = cache.Getter.GetVaccinations(ctx)
		if err != nil {
			entry.once = nil
		}
		cache.setCacheEntry("vaccinations", entry)
	})

	return entry.entries, err
}

// GetCases retrieves all COVID-19 cases.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetCases(ctx context.Context) (results []Measurement, err error) {
	entry := cache.getCacheEntry("cases")
	entry.once.Do(func() {
		entry.entries, err = cache.Getter.GetCases(ctx)
		if err != nil {
			entry.once = nil
		}
		cache.setCacheEntry("cases", entry)
	})

	return entry.entries, err
}

// GetMortality retrieves all COVID-19 deaths.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetMortality(ctx context.Context) (results []Measurement, err error) {
	entry := cache.getCacheEntry("mortality")
	entry.once.Do(func() {
		entry.entries, err = cache.Getter.GetMortality(ctx)
		if err != nil {
			entry.once = nil
		}
		cache.setCacheEntry("mortality", entry)
	})

	return entry.entries, err
}

// GetHospitalisations retrieves all hospitalisations.  If a valid cached result exists, that is returned instead.
func (cache *Cache) GetHospitalisations(ctx context.Context) (results []Measurement, err error) {
	entry := cache.getCacheEntry("hospitalisations")
	entry.once.Do(func() {
		entry.entries, err = cache.Getter.GetHospitalisations(ctx)
		if err != nil {
			entry.once = nil
		}
		cache.setCacheEntry("hospitalisations", entry)
	})

	return entry.entries, err
}
