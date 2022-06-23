package mortality

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/clambin/sciensano/apiclient/sciensano"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	table2 "github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APIClient   fetcher.Fetcher
}

// GroupedMortalityEntry contains all the values for the (grouped) mortality figures
type GroupedMortalityEntry struct {
	Name   string
	Values []int
}

// Get returns all mortality figures
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Mortality", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeMortality); err2 == nil {
			output = table2.NewFromAPIResponse(apiResult)
		}
		return
	})
}

// GetByRegion returns all mortality figures, grouped by region
func (r *Reporter) GetByRegion() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("MortalityByRegion", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeMortality); err2 == nil {
			output = table2.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		}
		return
	})
}

// GetByAgeGroup returns all Mortality, grouped by age group
func (r *Reporter) GetByAgeGroup() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("MortalityByAgeGroup", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeMortality); err2 == nil {
			output = table2.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByAgeGroup)
		}
		return
	})
}
