package fetcher

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

type cache struct {
	fetcher      Fetcher
	dataType     int
	entries      []apiclient.APIResponse
	lastModified time.Time
	lastChecked  time.Time
	expiry       time.Duration
	lock         sync.RWMutex
}

func (c *cache) Run(ctx context.Context, interval time.Duration) {
	if err := c.refresh(ctx); err != nil {
		log.WithError(err).Errorf("failed to get data for %s", c.fetcher.DataTypes()[c.dataType])
	}

	// add some jitter so all caches don't refresh at the same time
	jitter := rand.Int63n(60) - 30
	ticker := time.NewTicker(interval + time.Duration(jitter)*time.Second)
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

func (c *cache) refresh(ctx context.Context) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var serverTimestamp time.Time

	if time.Since(c.lastChecked) > c.expiry {
		serverTimestamp, err = c.fetcher.GetLastUpdated(ctx, c.dataType)
		if err != nil {
			return err
		}
		c.lastChecked = time.Now()
	}

	if serverTimestamp.After(c.lastModified) {
		var entries []apiclient.APIResponse
		if entries, err = c.fetcher.Fetch(ctx, c.dataType); err != nil {
			c.lastChecked = time.Time{}
			return err
		}
		c.entries = entries
		c.lastModified = serverTimestamp
	}

	return err
}
