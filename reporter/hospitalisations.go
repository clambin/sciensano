package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// HospitalisationsGetter contains all methods providing COVID-19-related hospitalisation figures
type HospitalisationsGetter interface {
	GetHospitalisations() (results *datasets.Dataset, err error)
	GetHospitalisationsByRegion() (results *datasets.Dataset, err error)
	GetHospitalisationsByProvince() (results *datasets.Dataset, err error)
}

// GetHospitalisations returns all hospitalisations
func (client *Client) GetHospitalisations() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Hospitalisations", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Hospitalisations"); found {
			output = datasets.GroupMeasurements(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}

// GetHospitalisationsByRegion returns all hospitalisations, grouped by region
func (client *Client) GetHospitalisationsByRegion() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("HospitalisationsByRegion", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Hospitalisations"); found {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}

// GetHospitalisationsByProvince returns all hospitalisations, grouped by province
func (client *Client) GetHospitalisationsByProvince() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("HospitalisationsByProvince", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Hospitalisations"); found {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByProvince)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}
