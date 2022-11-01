package fetcher

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type cache struct {
	fetcher   Fetcher
	dataType  int
	entries   []apiclient.APIResponse
	timestamp time.Time
	expiry    time.Duration
	lock      sync.RWMutex
}

func (c *cache) Run(ctx context.Context, interval time.Duration) {
	if err := c.refresh(ctx); err != nil {
		log.WithError(err).Errorf("failed to get data for %s", c.fetcher.DataTypes()[c.dataType])
	}

	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			if err := c.refresh(ctx); err != nil {
				log.WithError(err).Errorf("failed to get data for %s", c.fetcher.DataTypes()[c.dataType])
			}
		}
	}
	ticker.Stop()
}

func (c *cache) Get() []apiclient.APIResponse {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.entries
}

func (c *cache) refresh(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.timestamp.IsZero() {
		if time.Since(c.timestamp) < c.expiry {
			return nil
		}
		lastModified, err := c.fetcher.GetLastUpdated(ctx, c.dataType)
		if err != nil {
			return err
		}
		if !lastModified.After(c.timestamp) {
			c.timestamp = time.Now()
			return nil
		}
	}

	entries, err := c.fetcher.Fetch(ctx, c.dataType)
	if err == nil {
		c.entries = entries
		c.timestamp = time.Now()
	}
	return err
}
