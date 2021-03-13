package sciensano

import (
	"net/http"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	apiClient http.Client

	CacheDuration           time.Duration
	testCacheExpiry         time.Time
	testCache               []apiTestResponse
	vaccinationsCacheExpiry time.Time
	vaccinationsCache       []apiVaccinationsResponse
}

const baseURL = "https://epistat.sciensano.be/Data/"

type API interface {
	GetTests(end time.Time) (results []Test, err error)
	GetVaccinations(end time.Time) (results []Vaccination, err error)
	GetVaccinationsByAge(end time.Time, group string) (results []Vaccination, err error)
	GetVaccinationsByRegion(end time.Time, group string) (results []Vaccination, err error)
}
