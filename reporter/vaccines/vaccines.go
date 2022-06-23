package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/clambin/sciensano/apiclient/vaccines"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APIClient   fetcher.Fetcher
}

// Get returns all vaccines data
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Vaccines", func() (*data.Table, error) {
		apiResult, err2 := r.APIClient.Fetch(context.Background(), vaccines.TypeBatches)

		if err2 != nil {
			return nil, fmt.Errorf("call failed: %w", err2)
		}
		return table.NewFromAPIResponse(apiResult), nil

	})
}

// GetByManufacturer returns all hospitalisations, grouped by manufacturer
func (r *Reporter) GetByManufacturer() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("VaccinesByManufacturer", func() (*data.Table, error) {
		apiResult, err2 := r.APIClient.Fetch(context.Background(), vaccines.TypeBatches)

		if err2 != nil {
			return nil, fmt.Errorf("call failed: %w", err2)
		}
		return table.NewGroupedFromAPIResponse(apiResult, apiclient.GroupByManufacturer), nil
	})
}
