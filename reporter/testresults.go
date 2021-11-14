package reporter

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// TestResultsGetter contains all required methods to retrieve COVID-19 test results
type TestResultsGetter interface {
	GetTestResults(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetTestResults returns all COVID-19 test results up to endTime
func (client *Client) GetTestResults(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("TestResults", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetTestResults(ctx); err2 == nil {
			output = datasets.GroupMeasurements(apiResult)
		}
		return
	})
}
