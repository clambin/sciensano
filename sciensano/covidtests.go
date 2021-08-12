package sciensano

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"time"
)

type Test struct {
	Timestamp time.Time
	Total     int
	Positive  int
}

func (client *Client) GetTests(ctx context.Context, end time.Time) (results []Test, err error) {
	var apiResult []apiTestResponse

	if apiResult, err = client.getTests(ctx); err == nil {
		results = groupTests(apiResult, end)
	}

	return
}

type apiTestResponse struct {
	TimeStamp string `json:"DATE"`
	Province  string `json:"PROVINCE"`
	Region    string `json:"REGION"`
	Total     int    `json:"TESTS_ALL"`
	Positive  int    `json:"TESTS_ALL_POS"`
}

func (client *Client) getTests(ctx context.Context) (response []apiTestResponse, err error) {
	client.testsLock.Lock()
	defer client.testsLock.Unlock()

	client.init()

	if client.testsCache == nil || time.Now().After(client.testsCacheExpiry) {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+"/Data/COVID19BE_tests.json", nil)

		var resp *http.Response
		resp, err = client.HTTPClient.Do(req)

		if err == nil {
			if resp.StatusCode == http.StatusOK {
				var body []byte
				body, err = io.ReadAll(resp.Body)

				if err == nil {
					var stats []apiTestResponse
					err = json.Unmarshal(body, &stats)

					if err == nil {
						client.testsCache = stats
						client.testsCacheExpiry = time.Now().Add(client.CacheDuration)
					}
				}
			} else {
				err = errors.New(resp.Status)
			}
			_ = resp.Body.Close()
		}
	}
	return client.testsCache, err
}

func groupTests(apiResult []apiTestResponse, end time.Time) (results []Test) {
	// Store the totals in a map
	totals := make(map[time.Time]Test, 0)
	for _, entry := range apiResult {
		if ts, err2 := time.Parse("2006-01-02", entry.TimeStamp); err2 == nil {
			// Skip anything after the specified end date
			if ts.After(end) {
				continue
			}

			var current Test
			var ok bool
			if current, ok = totals[ts]; ok == false {
				current.Timestamp = ts
			}
			current.Total += entry.Total
			current.Positive += entry.Positive
			totals[ts] = current
		} else {
			log.WithFields(log.Fields{"err": err2, "timestamp": entry.TimeStamp}).Warning("could not parse timestamp from API. skipping entry")
		}
	}
	// For each entry in the map, create an entry in the results slice
	for _, entry := range totals {
		results = append(results, entry)
	}
	// Maps are iterated in random order. Sort the final slice
	sort.Slice(results, func(i, j int) bool { return results[i].Timestamp.Before(results[j].Timestamp) })

	return
}
