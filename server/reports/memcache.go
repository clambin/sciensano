package reports

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/clambin/go-common/tabulator"
	"time"
)

var _ Cache = &MemCache{}

type MemCache struct {
	cache      MemCacheClient
	expiration time.Duration
}

//go:generate mockery --name MemCacheClient
type MemCacheClient interface {
	Get(string) (*memcache.Item, error)
	Set(*memcache.Item) error
}

func NewMemCache(m MemCacheClient, expiration time.Duration) *MemCache {
	return &MemCache{cache: m, expiration: expiration}
}

func (m *MemCache) Get(key string) (*tabulator.Tabulator, error) {
	item, err := m.cache.Get(key)
	if err == nil {
		var report tabulator.Tabulator
		err = json.Unmarshal(item.Value, &report)
		return &report, err
	}

	if errors.Is(err, memcache.ErrCacheMiss) {
		err = ErrMissedCache
	}
	return nil, err
}

func (m *MemCache) Set(key string, table *tabulator.Tabulator) error {
	bytes, err := json.Marshal(table)
	if err == nil {
		err = m.cache.Set(&memcache.Item{
			Key:        key,
			Value:      bytes,
			Expiration: int32(m.expiration.Seconds()),
		})
	}
	if err != nil {
		err = fmt.Errorf("memcache set: %w", err)
	}
	return err
}
