package reports

import (
	bytes2 "bytes"
	"encoding/gob"
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
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			err = ErrMissedCache
		}
		return nil, err
	}

	var report tabulator.Tabulator
	err = gob.NewDecoder(bytes2.NewBuffer(item.Value)).Decode(&report)
	return &report, err
}

func (m *MemCache) Set(key string, table *tabulator.Tabulator) error {
	var body bytes2.Buffer
	if err := gob.NewEncoder(&body).Encode(table); err != nil {
		return fmt.Errorf("memcache set: encode: %w", err)
	}
	if err := m.cache.Set(&memcache.Item{
		Key:        key,
		Value:      body.Bytes(),
		Expiration: int32(m.expiration.Seconds()),
	}); err != nil {
		return fmt.Errorf("memcache set: %w", err)
	}

	return nil
}
