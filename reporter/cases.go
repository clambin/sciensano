package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases() (results *datasets.Dataset, err error)
	GetCasesByRegion() (results *datasets.Dataset, err error)
	GetCasesByProvince() (results *datasets.Dataset, err error)
	GetCasesByAgeGroup() (results *datasets.Dataset, err error)
}

// GetCases returns all cases
func (client *Client) GetCases() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Cases", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Cases"); found {
			output = datasets.GroupMeasurements(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetCasesByRegion returns all cases, grouped by region
func (client *Client) GetCasesByRegion() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("CasesByRegion", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Cases"); found {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetCasesByProvince returns all cases, grouped by province
func (client *Client) GetCasesByProvince() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("CasesByProvince", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Cases"); found {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByProvince)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetCasesByAgeGroup returns all cases, grouped by province
func (client *Client) GetCasesByAgeGroup() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("CasesByAgeGroup", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Cases"); found {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByAgeGroup)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}
