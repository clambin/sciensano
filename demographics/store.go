package demographics

import (
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
}

// DefaultAgeBrackets specifies the default age brackets for the demographics by age
var DefaultAgeBrackets = []float64{12, 16, 18, 25, 35, 45, 55, 65, 75, 85}

// Store holds the demographics data
type Store struct {
	// Retention specifies how long to cache the data
	Retention time.Duration
	// AgeBrackets specifies the age brackets to group the data in. Defaults to DefaultAgeBrackets
	AgeBrackets []float64
	// TempDirectory specifies the directory to use for temporary files. Uses system-specified tempdir is left blank
	TempDirectory string
	// URL is the URL that will be used to retrieve the data. Used for unit testing
	URL      string
	byAge    map[Bracket]int
	byRegion map[string]int
	lock     sync.RWMutex
	once     *sync.Once
	expiry   time.Time
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
	store.lock.Lock()
	defer store.lock.Unlock()

	byAge, _, err := store.update()
	if err == nil {
		count, ok = byAge[bracket]
	}
	return
}

// GetAgeBrackets returns all age brackets found in the demographics data
func (store *Store) GetAgeBrackets() (brackets []Bracket) {
	store.lock.Lock()
	defer store.lock.Unlock()

	byAge, _, err := store.update()
	if err == nil {
		for bracket := range byAge {
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
	store.lock.Lock()
	defer store.lock.Unlock()

	_, byRegion, err := store.update()
	if err == nil {
		count, ok = byRegion[region]
	}
	return
}

// GetRegions returns all regions found in the demographics data
func (store *Store) GetRegions() (regions []string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	_, byRegion, err := store.update()
	if err == nil {
		for region := range byRegion {
			regions = append(regions, region)
		}
	}
	return
}
