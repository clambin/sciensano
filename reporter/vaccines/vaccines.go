package vaccines

import (
	"errors"
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

// Get returns all vaccines data
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Vaccines", func() (*data.Table, error) {
		apiResult, ok := r.APICache.Get("Vaccines")

		if !ok {
			return nil, errors.New("cache does not contain Vaccines entries")
		}
		return table.NewFromAPIResponse(apiResult), nil

	})
}

// GetByManufacturer returns all hospitalisations, grouped by manufacturer
func (r *Reporter) GetByManufacturer() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("VaccinesByManufacturer", func() (*data.Table, error) {
		apiResult, ok := r.APICache.Get("Vaccines")

		if !ok {
			return nil, errors.New("cache does not contain Vaccines entries")
		}
		return table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByManufacturer), nil
	})
}
