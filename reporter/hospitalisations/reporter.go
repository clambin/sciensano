package hospitalisations

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

// Get returns all hospitalisations
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Hospitalisations", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Hospitalisations"); found {
			output = table.NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}

// GetByRegion returns all hospitalisations, grouped by region
func (r *Reporter) GetByRegion() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("HospitalisationsByRegion", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Hospitalisations"); found {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}

// GetByProvince returns all hospitalisations, grouped by province
func (r *Reporter) GetByProvince() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("HospitalisationsByProvince", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Hospitalisations"); found {
			output = table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByProvince)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}
