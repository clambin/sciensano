package sciensano

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

type TestResult struct {
	Timestamp time.Time
	Total     int
	Positive  int
}

func (client *Client) GetTests(ctx context.Context, end time.Time) (results []TestResult, err error) {
	var apiResult []apiTestResponse

	if apiResult, err = client.testResultsCache.GetTestResults(ctx); err == nil {
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

func groupTests(apiResult []apiTestResponse, end time.Time) (results []TestResult) {
	// Store the totals in a map
	totals := make(map[time.Time]TestResult, 0)
	for _, entry := range apiResult {
		if ts, err2 := time.Parse("2006-01-02", entry.TimeStamp); err2 == nil {
			// Skip anything after the specified end date
			if ts.After(end) {
				continue
			}

			var current TestResult
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
