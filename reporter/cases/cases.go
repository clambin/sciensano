package cases

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/clambin/sciensano/apiclient/sciensano"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APIClient   fetcher.Fetcher
}

// Get returns all cases
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Cases", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeCases); err2 == nil {
			output = table.NewFromAPIResponse(apiResult)
		}
		return
	})
}

// GetByRegion returns all cases, grouped by region
func (r *Reporter) GetByRegion() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("CasesByRegion", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeCases); err2 == nil {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		}
		return
	})
}

// GetByProvince returns all cases, grouped by province
func (r *Reporter) GetByProvince() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("CasesByProvince", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeCases); err2 == nil {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByProvince)
		}
		return
	})
}

// GetByAgeGroup returns all cases, grouped by province
func (r *Reporter) GetByAgeGroup() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("CasesByAgeGroup", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeCases); err2 == nil {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByAgeGroup)
		}
		return
	})
}
