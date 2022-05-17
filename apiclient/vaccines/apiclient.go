package vaccines

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/clambin/go-metrics/caller"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"time"
)

// Getter interface retrieves vaccine batches
//go:generate mockery --name Getter
type Getter interface {
	GetBatches(ctx context.Context) (batches []apiclient.APIResponse, err error)
}

var _ cache.Fetcher = &Client{}
var _ Getter = &Client{}

// Client calls the API to retrieve vaccine batches
type Client struct {
	caller.Caller
	URL string
	cache.Cache
}

const baseURL = "https://covid-vaccinatie.be"

func (client *Client) getURL() (url string) {
	url = baseURL
	if client.URL != "" {
		url = client.URL
	}
	return
}

// APIBatchResponse represents one batch of vaccines
type APIBatchResponse struct {
	Date         Timestamp `json:"date"`
	Manufacturer string    `json:"manufacturer"`
	Amount       int       `json:"amount"`
}

var _ apiclient.APIResponse = &APIBatchResponse{}

// GetTimestamp returns the batch's timestamp
func (b APIBatchResponse) GetTimestamp() time.Time {
	return b.Date.Time
}

// GetGroupFieldValue returns the value of a groupable field.  Not used for APIBatchResponse.
func (b APIBatchResponse) GetGroupFieldValue(groupField int) (value string) {
	if groupField == apiclient.GroupByManufacturer {
		value = b.Manufacturer
	}
	return
}

// GetTotalValue returns the entry's total number of vaccines
func (b APIBatchResponse) GetTotalValue() float64 {
	return float64(b.Amount)
}

// GetAttributeNames returns the names of the types of vaccinations
func (b APIBatchResponse) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (b APIBatchResponse) GetAttributeValues() (values []float64) {
	return []float64{float64(b.Amount)}
}

// Timestamp representation for APIBatchResponse. Needed to unmarshal the date as received from the API
type Timestamp struct {
	time.Time
}

// UnmarshalJSON unmarshals the Timestamp in a APIBatchResponse
func (date *Timestamp) UnmarshalJSON(b []byte) (err error) {
	var timestamp time.Time
	if timestamp, err = time.Parse(`"2006-01-02"`, string(b)); err == nil {
		date.Time = timestamp
	}
	return
}

// Update calls all endpoints and returns this to the caller. This is used by a cache to refresh its content
func (client *Client) Update(ctx context.Context, ch chan<- cache.FetcherResponse) {
	log.Debug("refreshing Vaccine API cache")
	start := time.Now()
	cache.Fetch(ctx, ch, "Vaccines", client.GetBatches)
	log.WithField("duration", time.Since(start)).Debugf("refreshed Vaccine API cache")
}

type apiBatchesResponse struct {
	Result struct {
		Delivered []*APIBatchResponse `json:"delivered"`
	} `json:"result"`
}

// GetBatches returns all vaccine batches
func (client *Client) GetBatches(ctx context.Context) (batches []apiclient.APIResponse, err error) {
	var stats apiBatchesResponse
	stats, err = client.call(ctx)
	if err != nil {
		return
	}

	batches = make([]apiclient.APIResponse, 0, len(stats.Result.Delivered))
	for _, entry := range stats.Result.Delivered {
		batches = append(batches, entry)
	}

	sort.Slice(batches, func(i, j int) bool { return batches[i].GetTimestamp().Before(batches[j].GetTimestamp()) })
	return
}

func (client *Client) call(ctx context.Context) (stats apiBatchesResponse, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.getURL()+"/api/v1/delivered.json", nil)

	var resp *http.Response
	resp, err = client.Caller.Do(req)

	if err != nil {
		return
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return stats, errors.New(resp.Status)
	}

	var body []byte
	if body, err = io.ReadAll(resp.Body); err == nil {
		err = json.Unmarshal(body, &stats)
	}
	return
}
