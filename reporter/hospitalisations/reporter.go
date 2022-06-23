package hospitalisations

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

// Get returns all hospitalisations
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Hospitalisations", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeHospitalisations); err2 == nil {
			output = table.NewFromAPIResponse(apiResult)
		}
		return
	})
}

// GetByRegion returns all hospitalisations, grouped by region
func (r *Reporter) GetByRegion() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("HospitalisationsByRegion", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeHospitalisations); err2 == nil {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		}
		return
	})
}

// GetByProvince returns all hospitalisations, grouped by province
func (r *Reporter) GetByProvince() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("HospitalisationsByProvince", func() (output *data.Table, err2 error) {
		var apiResult []apiclient.APIResponse
		if apiResult, err2 = r.APIClient.Fetch(context.Background(), sciensano.TypeHospitalisations); err2 == nil {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByProvince)
		}
		return
	})
}
