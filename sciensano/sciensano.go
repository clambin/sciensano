package sciensano

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	HTTPClient              *http.Client
	URL                     string
	CacheDuration           time.Duration
	testsCacheExpiry        time.Time
	testsCache              []apiTestResponse
	testsLock               sync.Mutex
	vaccinationsCacheExpiry time.Time
	vaccinationsCache       []apiVaccinationsResponse
	vaccinationsLock        sync.Mutex
}

const baseURL = "https://epistat.sciensano.be"

type API interface {
	GetTests(ctx context.Context, end time.Time) (results []Test, err error)
	GetVaccinations(ctx context.Context, end time.Time) (results []Vaccination, err error)
	GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
	GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
}

func (client *Client) init() {
	if client.URL == "" {
		client.URL = baseURL
	}
	if client.HTTPClient == nil {
		client.HTTPClient = &http.Client{}
	}
}
