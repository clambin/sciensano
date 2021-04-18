package sciensano

import (
	"net/http"
	"sync"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	HTTPClient              *http.Client
	CacheDuration           time.Duration
	testsCacheExpiry        time.Time
	testsCache              []apiTestResponse
	testsLock               sync.Mutex
	vaccinationsCacheExpiry time.Time
	vaccinationsCache       []apiVaccinationsResponse
	vaccinationsLock        sync.Mutex
}

const baseURL = "https://epistat.sciensano.be/Data/"

type API interface {
	GetTests(end time.Time) (results []Test, err error)
	GetVaccinations(end time.Time) (results []Vaccination, err error)
	GetVaccinationsByAge(end time.Time) (results map[string][]Vaccination, err error)
	GetVaccinationsByRegion(end time.Time) (results map[string][]Vaccination, err error)
}
