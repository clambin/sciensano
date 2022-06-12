package mortality

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	apiCache "github.com/clambin/sciensano/apiclient/cache"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	table2 "github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APICache    apiCache.Holder
}

// GroupedMortalityEntry contains all the values for the (grouped) mortality figures
type GroupedMortalityEntry struct {
	Name   string
	Values []int
}

// Get returns all mortality figures
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Mortality", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Mortality"); found {
			output = table2.NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Mortality entries")
		}
		return
	})
}

// GetByRegion returns all mortality figures, grouped by region
func (r *Reporter) GetByRegion() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("MortalityByRegion", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Mortality"); found {
			output = table2.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Mortality entries")
		}
		return
	})
}

// GetByAgeGroup returns all Mortality, grouped by age group
func (r *Reporter) GetByAgeGroup() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("MortalityByAgeGroup", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Mortality"); found {
			output = table2.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByAgeGroup)
		} else {
			err2 = fmt.Errorf("cache does not contain Mortality entries")
		}
		return
	})
}
