package reporter

import (
	"fmt"
	"github.com/clambin/simplejson/v3/dataset"
)

// TestResultsGetter contains all required methods to retrieve COVID-19 test results
type TestResultsGetter interface {
	GetTestResults() (results *dataset.Dataset, err error)
}

// GetTestResults returns all COVID-19 test results up to endTime
func (client *Client) GetTestResults() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("TestResults", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("TestResults"); found {
			output = NewFromAPIResponse(apiResult)
			output.AddColumn("rate", func(values map[string]float64) (rate float64) {
				if values["total"] != 0 {
					rate = values["positive"] / values["total"]
				}
				return
			})
		} else {
			err2 = fmt.Errorf("cache does not contain TestResults entries")
		}
		return
	})
}
