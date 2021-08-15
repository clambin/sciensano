package demographics

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// DefaultAgeBrackets specifies the default age brackets for the demographics by age
var DefaultAgeBrackets = []float64{12, 16, 18, 25, 35, 45, 55, 65, 75, 85}

// Server updates demographics on a recurring basis
type Server struct {
	// AgeBrackets configures in what age brackets the age demographics should be grouped. Defaults to DefaultAgeBrackets
	AgeBrackets []float64
	// TempDirectory specifies the directory to use for temporary files. Uses system-specified tempdir is left blank
	TempDirectory string
	// HTTPClient used to retrieve the data from the mock
	HTTPClient *http.Client
	// URL is the URL that will be used to retrieve the data. Used for unit testing
	URL      string
	byAge    map[Bracket]int
	byRegion map[string]int
	lock     sync.RWMutex
}

// New creates a new demographics mock
func New() *Server {
	return &Server{
		AgeBrackets: DefaultAgeBrackets,
		HTTPClient:  &http.Client{},
	}
}

// Run the mock, updating demographics at the specified interval
func (server *Server) Run(ctx context.Context, interval time.Duration) (err error) {
	err = server.update(ctx)

	ticker := time.NewTicker(interval)

	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			err = server.update(ctx)
		}
	}

	return
}

// AvailableData returns true if the mock has downloaded demographics data
func (server *Server) AvailableData() bool {
	server.lock.RLock()
	defer server.lock.RUnlock()

	return len(server.byRegion) > 0 || len(server.byAge) > 0
}

// GetByAge returns the number of people within the specified age bracket
func (server *Server) GetByAge(bracket Bracket) (count int, ok bool) {
	server.lock.RLock()
	defer server.lock.RUnlock()

	count, ok = server.byAge[bracket]
	return
}

// GetAgeBrackets returns all age brackets
func (server *Server) GetAgeBrackets() (brackets []Bracket) {
	server.lock.RLock()
	defer server.lock.RUnlock()

	for bracket := range server.byAge {
		brackets = append(brackets, bracket)
	}
	return
}

// GetByRegion returns the number of people within the specified region.
// Note: current data source doesn't differentiate for ostbelgien. This is currently handled by hard-coding that data.
func (server *Server) GetByRegion(region string) (count int, ok bool) {
	server.lock.RLock()
	defer server.lock.RUnlock()

	count, ok = server.byRegion[region]
	return
}

// GetRegions returns all regions
func (server *Server) GetRegions() (regions []string) {
	server.lock.RLock()
	defer server.lock.RUnlock()

	for region := range server.byRegion {
		regions = append(regions, region)
	}
	return
}
