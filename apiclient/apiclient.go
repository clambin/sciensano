package apiclient

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// Getter interface exposes the different supported Sciensano APIs
//go:generate mockery --name Getter
type Getter interface {
	GetTestResults(ctx context.Context) (results []*APITestResultsResponse, err error)
	GetVaccinations(ctx context.Context) (results []*APIVaccinationsResponse, err error)
	GetCases(ctx context.Context) (results []*APICasesResponse, err error)
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
func (client *Client) call(ctx context.Context, endpoint string, results interface{}) (err error) {
	target := client.getURL() + "/Data/" + endpoint

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)

	var resp *http.Response
	resp, err = client.HTTPClient.Do(req)

	if err != nil {
		return
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	var body []byte
	body, _ = io.ReadAll(resp.Body)
	err = json.Unmarshal(body, results)

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
