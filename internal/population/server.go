package population

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/v2/internal/population/bracket"
	"log/slog"
	"math"
	"sync"
	"time"
)

// Server imports the demographics data on a regular basis and exposes data APIs to callers
type Server struct {
	Waiter
	Path     string
	Interval time.Duration
	Logger   *slog.Logger
	mtime    time.Time
	byRegion map[string]int
	byAge    map[int]int
	lock     sync.RWMutex
}

// Run imports the latest demographics data on a regular basis
func (s *Server) Run(ctx context.Context) error {
	if err := s.update(); err != nil {
		return fmt.Errorf("population load failed: %w", err)
	}
	s.Ready()

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.update(); err != nil {
				s.Logger.Error("failed to read demographics file", "err", err)
			}
		}
	}
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
