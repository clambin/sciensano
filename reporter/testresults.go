package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/reporter/datasets"
)

// TestResultsGetter contains all required methods to retrieve COVID-19 test results
type TestResultsGetter interface {
	GetTestResults() (results *datasets.Dataset, err error)
}

// GetTestResults returns all COVID-19 test results up to endTime
func (client *Client) GetTestResults() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("TestResults", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("TestResults"); found {
			output = datasets.GroupMeasurements(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain TestResults entries")
		}
		return
	})
}
