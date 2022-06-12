package sciensano

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Getter interface exposes the different supported Sciensano APIs
//go:generate mockery --name Getter
type Getter interface {
	GetTestResults(ctx context.Context) (results []apiclient.APIResponse, err error)
	GetVaccinations(ctx context.Context) (results []apiclient.APIResponse, err error)
	GetCases(ctx context.Context) (results []apiclient.APIResponse, err error)
	GetMortality(ctx context.Context) (results []apiclient.APIResponse, err error)
	GetHospitalisations(ctx context.Context) (results []apiclient.APIResponse, err error)
}

var _ Getter = &Client{}
var _ cache.Fetcher = &Client{}

// Client calls the different Reporter APIs
type Client struct {
	client.Caller
	URL string
}

// Update calls all endpoints and returns this to the caller. This is used by a cache to refresh its content
func (client *Client) Fetch(ctx context.Context, ch chan<- cache.FetcherResponse) {
	start := time.Now()
	log.Debug("refreshing Reporter API cache")

	getters := map[string]func(context.Context) ([]apiclient.APIResponse, error){
		"Vaccinations":     client.GetVaccinations,
		"TestResults":      client.GetTestResults,
		"Hospitalisations": client.GetHospitalisations,
		"Mortality":        client.GetMortality,
		"Cases":            client.GetCases,
	}

	const maxParallel = 3
	s := semaphore.NewWeighted(maxParallel)

	for name, getter := range getters {
		_ = s.Acquire(ctx, 1)
		go func(name string, getter func(context.Context) ([]apiclient.APIResponse, error)) {
			cache.Fetch(ctx, ch, name, getter)
			s.Release(1)
		}(name, getter)
	}

	_ = s.Acquire(ctx, maxParallel)

	log.WithField("duration", time.Since(start)).Debugf("refreshed Reporter API cache")
}

const baseURL = "https://epistat.sciensano.be"

func (client *Client) getURL() (url string) {
	url = baseURL
	if client.URL != "" {
		url = client.URL
	}
	return
}

var endpoints = map[string]string{
	"cases":            "COVID19BE_CASES_AGESEX.json",
	"hospitalisations": "COVID19BE_HOSP.json",
	"mortality":        "COVID19BE_MORT.json",
	"tests":            "COVID19BE_tests.json",
	"vaccinations":     "COVID19BE_VACC.json",
}

// call is a generic function to call the Reporter API endpoints
func (client *Client) call(ctx context.Context, category string) (body []byte, err error) {
	endpoint, ok := endpoints[category]
	if !ok {
		err = fmt.Errorf("unsupporter category: %s", category)
		return
	}

	target := client.getURL() + "/Data/" + endpoint

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)

	var resp *http.Response
	resp, err = client.Caller.Do(req)

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	return io.ReadAll(resp.Body)
}

// TimeStamp represents a timestamp in the API responder. Needed for parsing purposes
type TimeStamp struct {
	time.Time
}

// UnmarshalJSON unmarshalls a TimeStamp from the API responder.
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

// copyMaybeSort accomplishes two goals: it recasts the individual API response type to the generic apiclient.APIResponse.
// Secondly, it checks if the data is in order. If not, it sorts it.
func copyMaybeSort[T apiclient.APIResponse](input []T) (output []apiclient.APIResponse) {
	output = make([]apiclient.APIResponse, len(input))
	var timestamp time.Time
	var mustSort bool
	for idx, entry := range input {
		output[idx] = entry
		if entry.GetTimestamp().Before(timestamp) {
			mustSort = true
		}
		timestamp = entry.GetTimestamp()
	}
	if mustSort {
		sort.Slice(output, func(i, j int) bool { return output[i].GetTimestamp().Before(output[j].GetTimestamp()) })
	}
	return
}
