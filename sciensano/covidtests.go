package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"sort"
	"time"
)

// TestsGetter contains all required methods to retrieve COVID-19 test results
type TestsGetter interface {
	GetTests(ctx context.Context, end time.Time) (results []TestResult, err error)
}

// TestResult represents results of administered COVID-19 tests
type TestResult struct {
	// Timestamp is the day the tests were administered
	Timestamp time.Time
	// The Total number of tests administered
	Total int
	// The Positive number of tests
	Positive int
}

// GetTests returns all COVID-19 tests up to endTime
func (client *Client) GetTests(ctx context.Context, endTime time.Time) (results []TestResult, err error) {
	var apiResult []*apiclient.APITestResultsResponse

	if apiResult, err = client.APIClient.GetTestResults(ctx); err == nil {
		results = groupTests(apiResult, endTime)
	}

	return
}

func groupTests(apiResult []*apiclient.APITestResultsResponse, end time.Time) (results []TestResult) {
	// Store the totals in a map
	totals := make(map[time.Time]TestResult, 0)
	for _, entry := range apiResult {
		// Skip anything after the specified end date
		if entry.TimeStamp.Time.After(end) {
			continue
		}

		var current TestResult
		var ok bool
		if current, ok = totals[entry.TimeStamp.Time]; ok == false {
			current.Timestamp = entry.TimeStamp.Time
		}

		current.Total += entry.Total
		current.Positive += entry.Positive
		totals[entry.TimeStamp.Time] = current
	}
	// For each entry in the map, create an entry in the results slice
	for _, entry := range totals {
		results = append(results, entry)
	}
	// Maps are iterated in random order. Sort the final slice
	sort.Slice(results, func(i, j int) bool { return results[i].Timestamp.Before(results[j].Timestamp) })

	return
}
