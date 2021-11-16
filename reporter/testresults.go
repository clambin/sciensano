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
	return client.ReportCache.MaybeGenerate("TestResults", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("TestResults"); found {
			output = datasets.GroupMeasurements(apiResult)
			calculateRatio(output)
		} else {
			err2 = fmt.Errorf("cache does not contain TestResults entries")
		}
		return
	})
}

func calculateRatio(input *datasets.Dataset) {
	ratios := make([]float64, len(input.Timestamps))
	for index := range input.Timestamps {
		var ratio float64
		if input.Groups[0].Values[index] != 0 {
			ratio = input.Groups[1].Values[index] / input.Groups[0].Values[index]
		}
		ratios[index] = ratio
	}

	input.Groups = append(input.Groups, datasets.DatasetGroup{
		Name:   "rate",
		Values: ratios,
	})
}
