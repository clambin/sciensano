package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/reporter/datasets"
)

// MortalityGetter contains all methods providing COVID-19 mortality
type MortalityGetter interface {
	GetMortality() (results *datasets.Dataset, err error)
	GetMortalityByRegion() (results *datasets.Dataset, err error)
	GetMortalityByAgeGroup() (results *datasets.Dataset, err error)
}

// GroupedMortalityEntry contains all the values for the (grouped) mortality figures
type GroupedMortalityEntry struct {
	Name   string
	Values []int
}

// GetMortality returns all mortality figures
func (client *Client) GetMortality() (results *datasets.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("Mortality", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Mortality"); found {
			output = datasets.NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Mortality entries")
		}
		return
	})
}

// GetMortalityByRegion returns all mortality figures, grouped by region
func (client *Client) GetMortalityByRegion() (results *datasets.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("MortalityByRegion", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Mortality"); found {
			output = datasets.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Mortality entries")
		}
		return
	})
}

// GetMortalityByAgeGroup returns all Mortality, grouped by age group
func (client *Client) GetMortalityByAgeGroup() (results *datasets.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("MortalityByAgeGroup", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Mortality"); found {
			output = datasets.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByAgeGroup)
		} else {
			err2 = fmt.Errorf("cache does not contain Mortality entries")
		}
		return
	})
}
