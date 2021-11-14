package reporter

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByProvince(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetCases returns all cases
func (client *Client) GetCases(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Cases", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetCases(ctx); err2 == nil {
			output = datasets.GroupMeasurements(apiResult)
		}
		return
	})
}

// GetCasesByRegion returns all cases, grouped by region
func (client *Client) GetCasesByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("CasesByRegion", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetCases(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByRegion)
		}
		return
	})
}

// GetCasesByProvince returns all cases, grouped by province
func (client *Client) GetCasesByProvince(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("CasesByProvince", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetCases(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByProvince)
		}
		return
	})
}

// GetCasesByAgeGroup returns all cases, grouped by province
func (client *Client) GetCasesByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("CasesByAgeGroup", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetCases(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByAgeGroup)
		}
		return
	})
}
