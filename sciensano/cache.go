package sciensano

import (
	"github.com/clambin/sciensano/sciensano/datasets"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
	"time"
)

type Cache struct {
	Duration time.Duration
	lock     sync.Mutex
	entries  map[string]*CacheEntry
}

type CacheEntry struct {
	Once   sync.Once
	Data   *datasets.Dataset
	expiry time.Time
}

func NewCache(duration time.Duration) *Cache {
	return &Cache{
		Duration: duration,
		entries:  make(map[string]*CacheEntry),
	}
}

// Prometheus metrics
var (
	metricCacheMiss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_reports_cache_miss_total",
		Help: "Number of Sciensano reports not served from cache",
	}, []string{"endpoint"})
	metricCacheCall = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sciensano_reports_cache_total",
		Help: "Number of Sciensano reports attempted to be served cache",
	}, []string{"endpoint"})
)

func (cache *Cache) Load(name string) (entry *CacheEntry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	var ok bool
	if entry, ok = cache.entries[name]; ok == false {
		entry = &CacheEntry{}
	}
	if ok == false || time.Now().After(entry.expiry) {
		entry.Once = sync.Once{}
		entry.expiry = time.Now().Add(cache.Duration)
		// register it so only one call to Load will set up a new Once
		cache.entries[name] = entry
		metricCacheMiss.WithLabelValues(name).Add(1.0)
	}
	metricCacheCall.WithLabelValues(name).Add(1.0)
	return
}

func (cache *Cache) Save(name string, entry *CacheEntry) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	entry.expiry = time.Now().Add(cache.Duration)
	cache.entries[name] = entry
}

func (cache *Cache) Clear(name string) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	delete(cache.entries, name)
}
