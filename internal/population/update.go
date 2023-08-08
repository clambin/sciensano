package population

import (
	"golang.org/x/exp/slog"
	"os"
	"time"
)

func (s *Server) update() (err error) {
	var mtime time.Time
	var updated bool
	if mtime, updated, err = s.isUpdated(); err != nil || !updated {
		return
	}

	if err = s.process(); err == nil {
		s.mtime = mtime
	}
	return
}

func (s *Server) isUpdated() (mtime time.Time, updated bool, err error) {
	var stats os.FileInfo
	if stats, err = os.Stat(s.Path); err == nil {
		mtime = stats.ModTime()
		updated = mtime.After(s.mtime)
	}
	return
}

const ostbelgienPopulation = 78000

func (s *Server) process() error {
	slog.Info("loading demographics")
	start := time.Now()
	byRegion, byAge, err := groupPopulation(s.Path)
	if err != nil {
		return err
	}

	// demographic figures counts Ostbelgien as part of Wallonia. Hardcode the split here.
	// yes, it's ugly. :-)
	_, found := byRegion["Ostbelgien"]
	if !found {
		byRegion["Ostbelgien"] = ostbelgienPopulation
		population := byRegion["Wallonia"]
		population -= ostbelgienPopulation
		byRegion["Wallonia"] = population
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.byRegion = byRegion
	s.byAge = byAge

	slog.Info("loaded demographics", "duration", time.Since(start))
	return nil
}
