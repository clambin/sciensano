package reports

import (
	"errors"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

type Cache interface {
	Set(string, *tabulator.Tabulator) error
	Get(string) (*tabulator.Tabulator, error)
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

type ReportCache struct {
	Cache
	lock    sync.RWMutex
	keys    set.Set[string]
	keyLock keyLock
}

var ErrMissedCache = errors.New("not in cache")

// MaybeGenerate loads a report from LocalCache, or generates it if the report does not exist or is expired
func (c *ReportCache) MaybeGenerate(key string, generate func() (*tabulator.Tabulator, error)) (*tabulator.Tabulator, error) {
	c.keyLock.Lock(key)
	defer c.keyLock.Unlock(key)

	metricCacheCall.WithLabelValues(key).Add(1)

	report, err := c.Get(key)
	if err == nil || !errors.Is(err, ErrMissedCache) {
		return report, err
	}

	metricCacheMiss.WithLabelValues(key).Add(1.0)

	if report, err = generate(); err != nil {
		return nil, err
	}

	err = c.Set(key, report)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if c.keys == nil {
		c.keys = set.Create[string]()
	}
	c.keys.Add(key)

	return report, nil
}

func (c *ReportCache) Stats() map[string]int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	stats := make(map[string]int)
	for _, key := range c.keys.List() {
		table, err := c.Cache.Get(key)
		var count int
		if err == nil {
			count = table.Size()
		}
		stats[key] = count
	}
	return stats
}

type keyLock struct {
	keys map[string]*sync.Mutex
	lock sync.Mutex
}

func (kl *keyLock) Lock(key string) {
	kl.lock.Lock()
	if kl.keys == nil {
		kl.keys = make(map[string]*sync.Mutex)
	}
	if _, ok := kl.keys[key]; !ok {
		kl.keys[key] = new(sync.Mutex)
	}
	kl.lock.Unlock()
	kl.keys[key].Lock()
}

func (kl *keyLock) Unlock(key string) {
	kl.lock.Lock()
	defer kl.lock.Unlock()
	kl.keys[key].Unlock()
}
