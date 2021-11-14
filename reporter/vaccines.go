package reporter

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// VaccinesGetter contains all required methods to retrieve vaccines data
type VaccinesGetter interface {
	GetVaccines(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetVaccines returns all vaccines data
func (client *Client) GetVaccines(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Vaccines", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Vaccines.GetBatches(ctx); err2 == nil {
			output = datasets.GroupMeasurements(apiResult)
		}
		return
	})
}
