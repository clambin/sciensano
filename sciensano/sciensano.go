package sciensano

import (
	"context"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	testResultsCache  *TestResultsCache
	vaccinationsCache *VaccinationsCache
}

const baseURL = "https://epistat.sciensano.be"

type API interface {
	GetTests(ctx context.Context, end time.Time) (results []TestResult, err error)
	GetVaccinations(ctx context.Context, end time.Time) (results []Vaccination, err error)
	GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
	GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
}

// NewClient creates a new Client
func NewClient(duration time.Duration) *Client {
	return &Client{
		testResultsCache:  NewTestResultsCache(duration),
		vaccinationsCache: NewVaccinationsCache(duration),
	}
}

// SetURL overrides the mock's URL.
// Used for unit testing. Not thread-safe. Use with caution
func (client *Client) SetURL(url string) {
	client.testResultsCache.URL = url
	client.vaccinationsCache.URL = url
}
