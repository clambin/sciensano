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

const ostbelgienPopulation = 78000

// GetForRegion returns the number of people in each region
func (s *Server) GetForRegion(region string) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	switch region {
	case "Ostbelgien":
		// demographic figures counts Ostbelgien as part of Wallonia. Hardcode the split here.
		// yes, it's ugly. :-)
		return ostbelgienPopulation
	case "Wallonia":
		return s.byRegion["Waals Gewest"] - ostbelgienPopulation
	default:
		return s.byRegion[translateRegion(region)]
	}
}

var regionTranslationTable = map[string]string{
	"Flanders": "Vlaams Gewest",
	"Wallonia": "Waals Gewest",
	"Brussels": "Brussels Hoofdstedelijk Gewest",
}

func translateRegion(input string) string {
	if translated, ok := regionTranslationTable[input]; ok {
		return translated
	}
	return input
}

// GetForAgeBracket returns the number of people within a specific age bracket. Set High to math.Inf(+1)
// to return all people older than a given age
func (s *Server) GetForAgeBracket(arguments bracket.Bracket) int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if arguments.High == 0 {
		arguments.High = math.Inf(+1)
	}

	var total int
	for age, count := range s.byAge {
		if float64(age) >= arguments.Low && float64(age) <= arguments.High {
			total += count
		}
	}
	return total
}
