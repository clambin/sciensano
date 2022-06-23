package fetcher

import (
	"github.com/clambin/sciensano/apiclient"
	"sync"
	"time"
)

type cacheEntry struct {
	entries []apiclient.APIResponse
	updated time.Time
	checked time.Time
	once    *sync.Once
}

func (e cacheEntry) shouldPoll(grace, expiration time.Duration) bool {
	sinceChecked := time.Since(e.checked)
	sinceUpdated := time.Since(e.updated)
	return sinceChecked > grace || sinceUpdated > expiration
}
