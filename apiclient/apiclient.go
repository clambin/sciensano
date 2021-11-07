package apiclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Getter interface exposes the different supported Sciensano APIs
//go:generate mockery --name Getter
type Getter interface {
	GetTestResults(ctx context.Context) (results []Measurement, err error)
	GetVaccinations(ctx context.Context) (results []Measurement, err error)
	GetCases(ctx context.Context) (results []Measurement, err error)
	GetMortality(ctx context.Context) (results []Measurement, err error)
	GetHospitalisations(ctx context.Context) (results []Measurement, err error)
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

// call is a generic function to call the Sciensano API endpoints
func (client *Client) call(ctx context.Context, endpoint string) (response io.ReadCloser, err error) {
	target := client.getURL() + "/Data/" + endpoint

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)

	var resp *http.Response
	if resp, err = client.HTTPClient.Do(req); err == nil {
		if resp.StatusCode == http.StatusOK {
			response = resp.Body
		} else {
			err = errors.New(resp.Status)
		}
	}
	return
}

// TimeStamp represents a timestamp in the API response. Needed for parsing purposes
type TimeStamp struct {
	time.Time
}

// UnmarshalJSON unmarshalls a TimeStamp from the API response.
func (ts *TimeStamp) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if len(s) != 12 || s[0] != '"' && s[11] != '"' {
		return fmt.Errorf("invalid timestamp: %s", s)
	}
	var year, month, day int
	year, errYear := strconv.Atoi(s[1:5])
	month, errMonth := strconv.Atoi(s[6:8])
	day, errDay := strconv.Atoi(s[9:11])

	if errYear != nil || errMonth != nil || errDay != nil {
		return fmt.Errorf("invalid timestamp: %s", s)
	}
	ts.Time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return
}
