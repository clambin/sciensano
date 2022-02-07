package demographics

import (
	"context"
	"github.com/clambin/sciensano/demographics/bracket"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
)

// Fetcher interface for population data
//go:generate mockery --name Fetcher
type Fetcher interface {
	// GetByAgeBracket returns the demographics for the specified AgeBracket
	GetByAgeBracket(arguments bracket.Bracket) (count int)
	// GetByRegion returns the demographics grouped by region
	GetByRegion() (figures map[string]int)
	// Run updates demographics and refreshes them on a periodic basis
	Run(ctx context.Context)
}

// Server imports the demographics data on a regular basis and exposes data APIs to callers
type Server struct {
	Path     string
	Interval time.Duration
	mtime    time.Time
	byRegion map[string]int
	byAge    map[int]int
	lock     sync.RWMutex
}

var _ Fetcher = &Server{}

// Run imports the latest demographics data on a regular basis
func (s *Server) Run(ctx context.Context) {
	log.Debug("first load of demographics file")
	err := s.update()
	if err != nil {
		log.WithError(err).Fatal("failed to read demographics file")
	}
	log.Debug("first load of demographics file done")

	ticker := time.NewTicker(s.Interval)

	for running := true; running; {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			running = false
		case <-ticker.C:
			err = s.update()
			if err != nil {
				log.WithError(err).Error("failed to read demographics file")
			}
		}
	}
	ticker.Stop()
}

// GetByRegion returns the number of people in each region
func (s *Server) GetByRegion() (response map[string]int) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	response = make(map[string]int)
	for key, value := range s.byRegion {
		response[key] = value
	}
	return
}

// GetByAgeBracket returns the number of people within a specific age bracket. Set High to math.Inf(+1)
// to return all people older than a given age
func (s *Server) GetByAgeBracket(arguments bracket.Bracket) (response int) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if arguments.High == 0 {
		arguments.High = math.Inf(+1)
	}
	for age, count := range s.byAge {
		if float64(age) >= arguments.Low && float64(age) <= arguments.High {
			response += count
		}
	}
	return
}
