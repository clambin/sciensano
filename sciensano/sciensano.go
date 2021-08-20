package sciensano

import (
	"context"
	"github.com/clambin/sciensano/sciensano/apiclient"
	"net/http"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	apiclient.APIClient
}

type APIClient interface {
	GetTests(ctx context.Context, end time.Time) (results []TestResult, err error)
	GetVaccinations(ctx context.Context, end time.Time) (results []Vaccination, err error)
	GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
	GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
}

// NewClient creates a new Client
func NewClient(duration time.Duration) *Client {
	return &Client{
		APIClient: &apiclient.Cache{
			APIClient: &apiclient.Client{HTTPClient: &http.Client{}},
			Retention: duration,
		},
	}
}
