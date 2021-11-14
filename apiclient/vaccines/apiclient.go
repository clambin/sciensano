package vaccines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clambin/sciensano/measurement"
	"io"
	"net/http"
	"sort"
	"time"
)

// Getter interface retrieves vaccine batches
//go:generate mockery --name Getter
type Getter interface {
	GetBatches(ctx context.Context) (batches []measurement.Measurement, err error)
}

// Client calls the API to retrieve vaccine batches
type Client struct {
	URL        string
	HTTPClient *http.Client
	measurement.Cache
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

var _ measurement.Measurement = &Batch{}

// GetTimestamp returns the batch's timestamp
func (b Batch) GetTimestamp() time.Time {
	return b.Date.Time
}

// GetGroupFieldValue returns the value of a groupable field.  Not used for Batch.
func (b Batch) GetGroupFieldValue(_ int) string {
	return ""
}

// GetTotalValue returns the entry's total number of vaccines
func (b Batch) GetTotalValue() float64 {
	return float64(b.Amount)
}

// GetAttributeNames returns the names of the types of vaccinations
func (b Batch) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (b Batch) GetAttributeValues() (values []float64) {
	return []float64{float64(b.Amount)}
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
func (client *Client) GetBatches(ctx context.Context) (batches []measurement.Measurement, err error) {
	return client.Cache.Call(ctx, "vaccines", client.getBatches)
}

func (client *Client) getBatches(ctx context.Context) (batches []measurement.Measurement, err error) {
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

	batches = make([]measurement.Measurement, 0, len(stats.Result.Delivered))
	for _, entry := range stats.Result.Delivered {
		batches = append(batches, entry)
	}

	sort.Slice(batches, func(i, j int) bool { return batches[i].GetTimestamp().Before(batches[j].GetTimestamp()) })

	return
}
