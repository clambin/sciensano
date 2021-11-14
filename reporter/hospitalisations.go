package reporter

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// HospitalisationsGetter contains all methods providing COVID-19-related hospitalisation figures
type HospitalisationsGetter interface {
	GetHospitalisations(ctx context.Context) (results *datasets.Dataset, err error)
	GetHospitalisationsByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetHospitalisationsByProvince(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetHospitalisations returns all hospitalisations
func (client *Client) GetHospitalisations(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Hospitalisations", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetHospitalisations(ctx); err2 == nil {
			output = datasets.GroupMeasurements(apiResult)
		}
		return
	})
}

// GetHospitalisationsByRegion returns all hospitalisations, grouped by region
func (client *Client) GetHospitalisationsByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("HospitalisationsByRegion", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetHospitalisations(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByRegion)
		}
		return
	})
}

// GetHospitalisationsByProvince returns all hospitalisations, grouped by province
func (client *Client) GetHospitalisationsByProvince(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("HospitalisationsByProvince", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetHospitalisations(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByProvince)
		}
		return
	})
}
