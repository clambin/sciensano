package reporter

import (
	"errors"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

// VaccinesGetter contains all required methods to retrieve vaccines data
type VaccinesGetter interface {
	GetVaccines() (results *datasets.Dataset, err error)
	GetVaccinesByManufacturer() (results *datasets.Dataset, err error)
}

// GetVaccines returns all vaccines data
func (client *Client) GetVaccines() (results *datasets.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("Vaccines", func() (*datasets.Dataset, error) {
		apiResult, ok := client.APICache.Get("Vaccines")

		if ok == false {
			return nil, errors.New("cache does not contain Vaccines entries")
		}
		return datasets.GroupMeasurements(apiResult), nil

	})
}

// GetVaccinesByManufacturer returns all hospitalisations, grouped by manufacturer
func (client *Client) GetVaccinesByManufacturer() (results *datasets.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("VaccinesByManufacturer", func() (*datasets.Dataset, error) {
		apiResult, ok := client.APICache.Get("Vaccines")

		if ok == false {
			return nil, errors.New("cache does not contain Vaccines entries")
		}
		return datasets.GroupMeasurementsByType(apiResult, measurement.GroupByManufacturer), nil
	})
}
