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
		var apiResult []apiclient.Measurement
		if apiResult, err = client.Getter.GetTestResults(ctx); err == nil {
			entry.Data = groupMeasurements(apiResult, apiclient.GroupByNone, NewTestResult)
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

// TestResult represents results of administered COVID-19 tests
type TestResult struct {
	// The Total number of tests administered
	Total int
	// The Positive number of tests
	Positive int
}

func NewTestResult() GroupedEntry {
	return &TestResult{}
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

// Add adds the passed test result values to its own values
func (entry *TestResult) Add(input apiclient.Measurement) {
	entry.Total += input.(*apiclient.APITestResultsResponseEntry).Total
	entry.Positive += input.(*apiclient.APITestResultsResponseEntry).Positive
}
