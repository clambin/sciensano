package reporter

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// MortalityGetter contains all methods providing COVID-19 mortality
type MortalityGetter interface {
	GetMortality(ctx context.Context) (results *datasets.Dataset, err error)
	GetMortalityByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetMortalityByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
}

// GroupedMortalityEntry contains all the values for the (grouped) mortality figures
type GroupedMortalityEntry struct {
	Name   string
	Values []int
}

// GetMortality returns all mortality figures
func (client *Client) GetMortality(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Mortality", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetMortality(ctx); err2 == nil {
			output = datasets.GroupMeasurements(apiResult)
		}
		return
	})
}

// GetMortalityByRegion returns all mortality figures, grouped by region
func (client *Client) GetMortalityByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("MortalityByRegion", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetMortality(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByRegion)
		}
		return
	})
}

// GetMortalityByAgeGroup returns all Mortality, grouped by age group
func (client *Client) GetMortalityByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("MortalityByAgeGroup", func() (output *datasets.Dataset, err2 error) {
		var apiResult []measurement.Measurement
		if apiResult, err2 = client.Sciensano.GetMortality(ctx); err2 == nil {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByAgeGroup)
		}
		return
	})
}
