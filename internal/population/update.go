package population

import (
	"os"
	"time"
)

func (s *Server) update() error {
	mtime, updated, err := s.isUpdated()
	if err != nil || !updated {
		return err
	}

	if err = s.process(); err == nil {
		s.mtime = mtime
	}
	return err
}

func (s *Server) isUpdated() (time.Time, bool, error) {
	stats, err := os.Stat(s.Path)
	if err != nil {
		return time.Time{}, false, err
	}

	mtime := stats.ModTime()
	return mtime, mtime.After(s.mtime), nil
}

func (s *Server) process() error {
	s.Logger.Info("loading demographics")
	start := time.Now()
	byRegion, byAge, err := groupPopulation(s.Path)
	if err == nil {
		s.lock.Lock()
		defer s.lock.Unlock()
		s.byRegion = byRegion
		s.byAge = byAge

		s.Logger.Info("loaded demographics", "duration", time.Since(start))
	}
	return err
}
