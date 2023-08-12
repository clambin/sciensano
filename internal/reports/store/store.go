package store

import (
	"errors"
	"github.com/clambin/go-common/tabulator"
	"golang.org/x/exp/slog"
	"sync"
)

type Store struct {
	Logger  *slog.Logger
	reports map[string]*tabulator.Tabulator
	lock    sync.RWMutex
}

var ErrNotFound = errors.New("report not found")

func (s *Store) Put(key string, report *tabulator.Tabulator) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.reports == nil {
		s.reports = make(map[string]*tabulator.Tabulator)
	}
	s.reports[key] = report
	s.Logger.Debug("report stored", "name", key, "rows", report.Size())
}

func (s *Store) Get(key string) (*tabulator.Tabulator, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.reports == nil {
		return nil, ErrNotFound
	}
	report, ok := s.reports[key]
	if !ok {
		return nil, ErrNotFound
	}
	return report, nil
}

func (s *Store) Keys() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var keys []string
	for key := range s.reports {
		keys = append(keys, key)
	}
	return keys
}
