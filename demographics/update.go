package demographics

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func (s *Server) update() (err error) {
	var mtime time.Time
	var updated bool
	mtime, updated, err = s.isUpdated()
	if err != nil || !updated {
		return
	}

	err = s.process()
	if err != nil {
		return
	}

	s.mtime = mtime
	return
}

func (s *Server) isUpdated() (mtime time.Time, updated bool, err error) {
	var stats os.FileInfo
	stats, err = os.Stat(s.Path)
	if err != nil {
		return
	}

	mtime = stats.ModTime()
	updated = mtime.After(s.mtime)
	return
}

const ostbelgienPopulation = 78000

func (s *Server) process() (err error) {
	log.Info("loading demographics")
	var (
		byRegion map[string]int
		byAge    map[int]int
	)

	byRegion, byAge, err = groupPopulation(s.Path)
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

	log.Info("loaded demographics")
	return
}
