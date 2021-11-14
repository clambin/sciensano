package sciensano

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/sciensano/measurement"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Getter interface exposes the different supported Sciensano APIs
//go:generate mockery --name Getter
type Getter interface {
	GetTestResults(ctx context.Context) (results []measurement.Measurement, err error)
	GetVaccinations(ctx context.Context) (results []measurement.Measurement, err error)
	GetCases(ctx context.Context) (results []measurement.Measurement, err error)
	GetMortality(ctx context.Context) (results []measurement.Measurement, err error)
	GetHospitalisations(ctx context.Context) (results []measurement.Measurement, err error)
	AutoRefresh(ctx context.Context)
}

// Client calls the different Sciensano APIs
type Client struct {
	URL        string
	HTTPClient *http.Client
	measurement.Cache
}

// AutoRefresh refreshes the cache on a period basis
func (client *Client) AutoRefresh(ctx context.Context) {
	client.refresh(ctx)
	ticker := time.NewTicker(client.Cache.Retention + 5*time.Second)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-ticker.C:
			client.refresh(ctx)
		}
	}
}

func (client *Client) refresh(ctx context.Context) {
	before := time.Now()
	log.Debug("refreshing API cache")

	const maxParallel = 3
	s := semaphore.NewWeighted(maxParallel)

	for _, getter := range []func(context.Context) ([]measurement.Measurement, error){
		client.GetVaccinations,
		client.GetTestResults,
		client.GetHospitalisations,
		client.GetMortality,
		client.GetCases,
	} {
		_ = s.Acquire(ctx, 1)
		go func(g func(context.Context) ([]measurement.Measurement, error)) {
			if _, err := g(ctx); err != nil {
				log.WithError(err).Warning("failed to call Sciensano endpoint")
			}
			s.Release(1)
		}(getter)
	}

	_ = s.Acquire(ctx, maxParallel)
	log.WithFields(log.Fields{
		"duration":  time.Now().Sub(before),
		"size":      client.CacheSize(),
		"retention": client.Retention,
	}).Debugf("refreshed API cache")
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
