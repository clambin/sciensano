package cache

import (
	"context"
	log "github.com/sirupsen/logrus"
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

func (s *cacher[T]) Get() T {
	_ = s.refresh(context.Background())
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.entries
}

func (s *cacher[T]) AutoRefresh(ctx context.Context, interval time.Duration) {
	err := s.refresh(ctx)
	if err != nil {
		log.WithError(err).WithField("target", s.GetTarget()).Error("failed to update cache")
	}

	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			err = s.refresh(ctx)
			if err != nil {
				log.WithError(err).WithField("target", s.GetTarget()).Error("failed to update cache")
			}
		}
	}
	ticker.Stop()
}

func (s *cacher[T]) refresh(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var serverTimestamp time.Time

	if time.Since(s.lastChecked) > s.expiry {
		var err error
		if serverTimestamp, err = s.GetLastModified(ctx); err != nil {
			return err
		}
		s.lastChecked = time.Now()
	}

	if serverTimestamp.After(s.lastModified) {
		entries, err := s.Fetch(ctx)
		if err != nil {
			s.lastChecked = time.Time{}
			return err
		}
		s.entries = entries
		s.lastModified = serverTimestamp
	}

	return nil
}
