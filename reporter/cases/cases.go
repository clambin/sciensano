package cases

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	apiCache "github.com/clambin/sciensano/apiclient/cache"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APICache    apiCache.Holder
}

// Get returns all cases
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Cases", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Cases"); found {
			output = table.NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetByRegion returns all cases, grouped by region
func (r *Reporter) GetByRegion() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("CasesByRegion", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Cases"); found {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetByProvince returns all cases, grouped by province
func (r *Reporter) GetByProvince() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("CasesByProvince", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Cases"); found {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByProvince)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}

// GetByAgeGroup returns all cases, grouped by province
func (r *Reporter) GetByAgeGroup() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("CasesByAgeGroup", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Cases"); found {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByAgeGroup)
		} else {
			err2 = fmt.Errorf("cache does not contain Cases entries")
		}
		return
	})
}
