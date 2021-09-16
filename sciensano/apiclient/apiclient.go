package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// APIClient interface exposes the different supported Sciensano APIs
//go:generate mockery --name APIClient
type APIClient interface {
	GetTestResults(ctx context.Context) (results []*APITestResultsResponse, err error)
	GetVaccinations(ctx context.Context) (results []*APIVaccinationsResponse, err error)
}

// Client calls the different sciensano APIs
type Client struct {
	URL        string
	HTTPClient *http.Client
}

const baseURL = "https://epistat.sciensano.be"

func (client *Client) getURL() (url string) {
	url = baseURL
	if client.URL != "" {
		url = client.URL
	}
	return
}

// TimeStamp represents a timestamp in the API response. Needed for parsing purposes
type TimeStamp struct {
	time.Time
}

// UnmarshalJSON unmarshalls a TimeStamp from the API response.
func (ts *TimeStamp) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	err = json.Unmarshal(b, &v)

	if err == nil {
		ts.Time, err = time.Parse("2006-01-02", v.(string))
	}
	return
}
