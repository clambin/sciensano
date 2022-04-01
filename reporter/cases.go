package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/simplejson/v3/dataset"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases() (results *dataset.Dataset, err error)
	GetCasesByRegion() (results *dataset.Dataset, err error)
	GetCasesByProvince() (results *dataset.Dataset, err error)
	GetCasesByAgeGroup() (results *dataset.Dataset, err error)
}

// GetCases returns all cases
func (client *Client) GetCases() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("Cases", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Cases"); found {
			output = NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetCasesByRegion returns all cases, grouped by region
func (client *Client) GetCasesByRegion() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("CasesByRegion", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Cases"); found {
			output = NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetCasesByProvince returns all cases, grouped by province
func (client *Client) GetCasesByProvince() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("CasesByProvince", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Cases"); found {
			output = NewGroupedFromAPIResponse(apiResult, apiclient.GroupByProvince)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetCasesByAgeGroup returns all cases, grouped by province
func (client *Client) GetCasesByAgeGroup() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("CasesByAgeGroup", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Cases"); found {
			output = NewGroupedFromAPIResponse(apiResult, apiclient.GroupByAgeGroup)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}
