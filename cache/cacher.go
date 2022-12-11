package cache

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
	"sync"
	"time"
)

type cacher[T any] struct {
	Fetcher[T]
	lock         sync.RWMutex
	entries      T
	expiry       time.Duration
	lastChecked  time.Time
	lastModified time.Time
}

func (s *cacher[T]) Get(_ context.Context) T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.entries
}

func (s *cacher[T]) AutoRefresh(ctx context.Context, interval time.Duration) {
	s.refresh(ctx)

	ticker := time.NewTicker(jitter(interval))
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			s.refresh(ctx)
		}
	}
	ticker.Stop()
}

func jitter(interval time.Duration) time.Duration {
	// randomize the interval (somewhat) so all caches don't try to update at the same time
	const window = 0.02 // 1% on either side of the interval
	seconds := interval.Microseconds()
	delta := int64(float64(seconds) * window)
	lowMark := seconds - delta/2
	j := rand.Int63n(delta)
	return time.Duration(lowMark+j) * time.Microsecond
}

func (s *cacher[T]) refresh(ctx context.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var serverTimestamp time.Time
	var err error

	if time.Since(s.lastChecked) > s.expiry {
		if serverTimestamp, err = s.GetLastModified(ctx); err == nil {
			s.lastChecked = time.Now()
		}
	}

	if serverTimestamp.After(s.lastModified) {
		var entries T
		if entries, err = s.Fetch(ctx); err == nil {
			s.entries = entries
			s.lastModified = serverTimestamp
		}
	}

	if err != nil {
		s.lastChecked = time.Time{}
		log.WithError(err).WithField("target", s.GetTarget()).Error("failed to update cache")
	}
}
