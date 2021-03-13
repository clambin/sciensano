package sciensano

import (
	"net/http"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	apiClient http.Client

	VaccinationsCacheDuration time.Duration
	vaccinationsCacheExpiry   time.Time
	vaccinationsCache         []apiVaccinationsResponse
}

const baseURL = "https://epistat.sciensano.be/Data/"
