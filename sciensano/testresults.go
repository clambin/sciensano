package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// TestResultsGetter contains all required methods to retrieve COVID-19 test results
type TestResultsGetter interface {
	GetTestResults(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetTestResults returns all COVID-19 test results up to endTime
func (client *Client) GetTestResults(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getTestResults(ctx, "GetTestResults", "TestResults")
}

func (client *Client) getTestResults(ctx context.Context, name, cacheEntryName string) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult apiclient.APITestResultsResponse
		if apiResult, err = client.Getter.GetTestResults(ctx); err == nil {
			entry.Data = groupTestResults(apiResult)
			client.cache.Save(cacheEntryName, entry)
		} else {
			client.cache.Clear(cacheEntryName)
		}
	})
	if err == nil && entry.Data != nil {
		results = entry.Data.Copy()
	}
	return
}

func groupTestResults(apiResult apiclient.APITestResultsResponse) (results *datasets.Dataset) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("groupTestResults")

	results = &datasets.Dataset{
		Timestamps: make([]time.Time, 0),
		Groups: []datasets.GroupedDatasetEntry{
			{
				Name:   "test_results",
				Values: make([]datasets.Copyable, 0),
			},
		},
	}

	currentTimestamp := time.Time{}
	currentEntry := &TestResult{}

	for _, entry := range apiResult {
		if !currentTimestamp.IsZero() && !currentTimestamp.Equal(entry.TimeStamp.Time) {
			results.Timestamps = append(results.Timestamps, currentTimestamp)
			results.Groups[0].Values = append(results.Groups[0].Values, currentEntry)
			currentEntry = &TestResult{}
		}

		currentTimestamp = entry.TimeStamp.Time
		currentEntry.Total += entry.Total
		currentEntry.Positive += entry.Positive
	}

	if !currentTimestamp.IsZero() {
		results.Timestamps = append(results.Timestamps, currentTimestamp)
		results.Groups[0].Values = append(results.Groups[0].Values, currentEntry)

	}
	return
}

// TestResult represents results of administered COVID-19 tests
type TestResult struct {
	// The Total number of tests administered
	Total int
	// The Positive number of tests
	Positive int
}

// Copy makes a copy of a TestResult
func (entry *TestResult) Copy() datasets.Copyable {
	return &TestResult{
		Total:    entry.Total,
		Positive: entry.Positive,
	}
}

// Ratio returns the positive rate for the test result
func (entry TestResult) Ratio() float64 {
	return float64(entry.Positive) / float64(entry.Total)
}
