package vaccines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// APIClient interface retrieves vaccine batches
//go:generate mockery --name APIClient
type APIClient interface {
	GetBatches(ctx context.Context) (batches []*Batch, err error)
}

// Client calls the API to retrieve vaccine batches
type Client struct {
	URL        string
	HTTPClient *http.Client
}

const baseURL = "https://covid-vaccinatie.be"

func (client *Client) getURL() (url string) {
	url = baseURL
	if client.URL != "" {
		url = client.URL
	}
	return
}

// Batch represents one batch of vaccines
type Batch struct {
	// Date is the date the batch was received
	Date Time
	// Manufacturer string
	// Amount is the number of vaccines in the batch
	Amount int
}

// Time representation for Batch. Needed to unmarshal the date as received from the API
type Time struct {
	time.Time
}

// UnmarshalJSON unmarshals the Time in a Batch
func (date *Time) UnmarshalJSON(b []byte) (err error) {
	var timestamp time.Time
	if timestamp, err = time.Parse(`"2006-01-02"`, string(b)); err == nil {
		date.Time = timestamp
	}
	return
}

// GetBatches returns all vaccine batches
func (client *Client) GetBatches(ctx context.Context) (batches []*Batch, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.getURL()+"/api/v1/delivered.json", nil)

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
	body, err = io.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("unable to parse vaccines response: %s", err.Error())
		return
	}

	var stats struct {
		Result struct {
			Delivered []*Batch `json:"delivered"`
		} `json:"result"`
	}

	err = json.Unmarshal(body, &stats)

	if err != nil {
		err = fmt.Errorf("unable to parse vaccines response: %s", err.Error())
		return
	}

	batches = stats.Result.Delivered
	sort.Slice(batches, func(i, j int) bool { return batches[i].Date.Time.Before(batches[j].Date.Time) })

	return
}

// AccumulateBatches accumulates the number of vaccines in the Batch list
func AccumulateBatches(batches []*Batch) (accumulated []*Batch) {
	var total int
	for _, batch := range batches {
		total += batch.Amount
		accumulated = append(accumulated, &Batch{
			Date:   batch.Date,
			Amount: total,
		})
	}
	return
}
