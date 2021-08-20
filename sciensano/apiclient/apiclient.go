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
	SetURL(url string)
}

// Client calls the different sciensano APIs
type Client struct {
	URL        string
	HTTPClient *http.Client
}

const baseURL = "https://epistat.sciensano.be"

func (client *Client) SetURL(url string) {
	client.URL = url
}

func (client *Client) getURL() (url string) {
	url = baseURL
	if client.URL != "" {
		url = client.URL
	}
	return
}

type TimeStamp struct {
	time.Time
}

func (ts *TimeStamp) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	err = json.Unmarshal(b, &v)

	if err == nil {
		ts.Time, err = time.Parse("2006-01-02", v.(string))
	}
	return
}
