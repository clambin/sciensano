package reporter

import (
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/measurement"
	"net/http"
	"time"
)

// Client queries different Reporter APIs
type Client struct {
	Sciensano sciensano.Getter
	Vaccines  vaccines.Getter
	Cache
}

// Reporter exposes the supported Reporter APIs
type Reporter interface {
	TestResultsGetter
	VaccinationGetter
	CasesGetter
	MortalityGetter
	HospitalisationsGetter
	VaccinesGetter
}

var _ Reporter = &Client{}

// NewCachedClient creates a new Client which caches results for duration interval
func NewCachedClient(duration time.Duration) *Client {
	return &Client{
		Sciensano: &sciensano.Client{
			HTTPClient: &http.Client{},
			Cache:      measurement.Cache{Retention: duration + 5*time.Second},
		},
		Vaccines: &vaccines.Client{
			HTTPClient: &http.Client{},
			Cache:      measurement.Cache{Retention: duration + 5*time.Second},
		},
		Cache: *NewCache(duration),
	}
}
