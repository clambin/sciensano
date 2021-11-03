package vaccines

import (
	"context"
	"sync"
	"time"
)

// Cache is a drop-in replacement for APIClient. It will cache the retrieved batches for the Retention duration.
type Cache struct {
	APIClient
	Retention   time.Duration
	lock        sync.Mutex
	batchesOnce *sync.Once
	batches     []*Batch
	testExpiry  time.Time
}

// GetBatches returns all vaccine batches
func (cache *Cache) GetBatches(ctx context.Context) (batches []*Batch, err error) {
	cache.lock.Lock()
	batches = cache.batches
	if cache.batchesOnce == nil || time.Now().After(cache.testExpiry) {
		cache.batchesOnce = &sync.Once{}
		cache.testExpiry = time.Now().Add(cache.Retention)
	}
	cache.lock.Unlock()

	cache.batchesOnce.Do(func() {
		if batches, err = cache.APIClient.GetBatches(ctx); err == nil {
			cache.batches = batches
		}
	})

	return
}
