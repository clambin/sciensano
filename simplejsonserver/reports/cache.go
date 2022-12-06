package reports

import (
	"errors"
	"github.com/clambin/go-common/cache"
	"github.com/clambin/go-common/tabulator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Cache caches the different reports generated by the server
type Cache struct {
	cache.Cacher[string, *Entry]
	lock sync.RWMutex
}

// Entry represents one report in the cache
type Entry struct {
	Once sync.Once
	Data *tabulator.Tabulator
}

// NewCache returns a new cache that caches the reports for the specified duration
func NewCache(duration time.Duration) *Cache {
	return &Cache{
		Cacher: cache.New[string, *Entry](duration, 5*time.Minute),
	}
}

// Prometheus metrics
var (
	metricCacheMiss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_reports_cache_miss_total",
		Help: "Number of Reporter reports not served from cache",
	}, []string{"endpoint"})
	metricCacheCall = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_reports_cache_total",
		Help: "Number of Reporter reports attempted to be served cache",
	}, []string{"endpoint"})
)

// Load returns a cached report. If the report does not exist, or is expired, it sets up Once to generate a new report
func (cache *Cache) Load(name string) (entry *Entry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	var ok bool
	if entry, ok = cache.Cacher.Get(name); !ok {
		entry = &Entry{}
	}
	if !ok || entry == nil {
		entry = &Entry{
			Once: sync.Once{},
		}
		// register it so only one call to Load will set up a new Once
		cache.Cacher.Add(name, entry)
		metricCacheMiss.WithLabelValues(name).Add(1.0)
	}
	metricCacheCall.WithLabelValues(name).Add(1.0)
	return
}

// Save stores a generated report in the cache
func (cache *Cache) Save(name string, entry *Entry) {
	cache.Cacher.Add(name, entry)
}

// MaybeGenerate loads a report from Cache, or generates it if the report does not exist or is expired
func (cache *Cache) MaybeGenerate(name string, generate func() (*tabulator.Tabulator, error)) (report *tabulator.Tabulator, err error) {
	entry := cache.Load(name)
	entry.Once.Do(func() {
		start := time.Now()
		log.WithField("report", name).Debug("generating report")
		entry.Data, err = generate()
		if err != nil {
			log.WithError(err).Warningf("failed to generate report '%s'", name)
			entry = nil
		}
		cache.Save(name, entry)
		log.WithFields(log.Fields{"report": name, "duration": time.Since(start)}).Debug("generated report")
	})
	if err == nil {
		report = entry.Data
		if report == nil {
			err = errors.New("not available")
		}
	}
	return
}

// Stats returns the stats of the cache
func (cache *Cache) Stats() (stats map[string]int) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	stats = make(map[string]int)
	for _, name := range cache.GetKeys() {
		entries, ok := cache.Get(name)
		var count int
		if ok && entries != nil && entries.Data != nil {
			count = len(entries.Data.GetTimestamps())
		}
		stats[name] = count
	}
	return
}
