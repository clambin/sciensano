package demographics

import (
	"context"
	"github.com/clambin/go-metrics"
	"sync"
	"time"
)

// Demographics interface giving access to available data
//go:generate mockery --name Demographics
type Demographics interface {
	// GetAgeGroupFigures returns the demographics grouped by age groups specified in AgeBrackets
	GetAgeGroupFigures() (figures map[string]int)
	// GetRegionFigures returns the demographics grouped by region
	GetRegionFigures() (figures map[string]int)
	// Stats returns statistics on the cache
	Stats() (stats map[string]int)
	// AutoRefresh periodically updates the cache
	AutoRefresh(ctx context.Context, interval time.Duration)
}

// DefaultAgeBrackets specifies the default age brackets for the demographics by age
var DefaultAgeBrackets = []float64{12, 16, 18, 25, 35, 45, 55, 65, 75, 85}

// Store holds the demographics data
type Store struct {
	AgeBrackets   []float64                // age brackets to group the data in. Defaults to DefaultAgeBrackets
	TempDirectory string                   // directory to use for temporary files. Uses system-specified tempdir if left blank
	URL           string                   // used to retrieve the data. Used for unit testing
	Metrics       metrics.APIClientMetrics // metrics to report API performance

	byAge    map[Bracket]int
	byRegion map[string]int
	lock     sync.RWMutex
	once     *sync.Once
	expiry   time.Time
}

// AutoRefresh periodically updates the cache
func (store *Store) AutoRefresh(ctx context.Context, interval time.Duration) {
	store.Update()

	ticker := time.NewTicker(interval)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			store.Update()
		}
	}
}

// GetAgeGroupFigures returns the demographics grouped by age groups specified in AgeBrackets
func (store *Store) GetAgeGroupFigures() (figures map[string]int) {
	figures = make(map[string]int)
	for _, ageGroup := range store.GetAgeBrackets() {
		figure, ok := store.GetByAge(ageGroup)
		if ok {
			figures[ageGroup.String()] = figure
		}
	}
	return
}

// GetByAge returns the total population in the specified age brackets
func (store *Store) GetByAge(bracket Bracket) (count int, ok bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	ok = false
	if store.byAge != nil {
		count, ok = store.byAge[bracket]
	}
	return
}

// GetAgeBrackets returns all age brackets found in the demographics data
func (store *Store) GetAgeBrackets() (brackets []Bracket) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	if store.byAge != nil {
		for bracket := range store.byAge {
			brackets = append(brackets, bracket)
		}
	}
	return
}

// GetRegionFigures returns the demographics grouped by region
func (store *Store) GetRegionFigures() (figures map[string]int) {
	figures = make(map[string]int)
	for _, region := range store.GetRegions() {
		figure, ok := store.GetByRegion(region)
		if ok {
			figures[region] = figure
		}
	}
	return
}

// GetByRegion returns the total population for the specified region
func (store *Store) GetByRegion(region string) (count int, ok bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	ok = false
	if store.byRegion != nil {
		count, ok = store.byRegion[region]
	}
	return
}

// GetRegions returns all regions found in the demographics data
func (store *Store) GetRegions() (regions []string) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	if store.byRegion != nil {
		for region := range store.byRegion {
			regions = append(regions, region)
		}
	}
	return
}

// Stats reports statistics on the case
func (store *Store) Stats() (stats map[string]int) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	stats = map[string]int{
		"Regions":     0,
		"AgeBrackets": 0,
	}

	if store.byAge != nil {
		stats["AgeBrackets"] = len(store.byAge)
	}
	if store.byRegion != nil {
		stats["Regions"] = len(store.byRegion)
	}
	return
}
