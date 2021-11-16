package reporter

import (
	"fmt"
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
	return client.ReportCache.MaybeGenerate("Vaccines", func() (output *datasets.Dataset, err2 error) {
		if apiResult, ok := client.APICache.Get("Vaccines"); ok {
			output = datasets.GroupMeasurements(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Vaccines entries")
		}
		return
	})
}

// GetVaccinesByManufacturer returns all hospitalisations, grouped by manufacturer
func (client *Client) GetVaccinesByManufacturer() (results *datasets.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("VaccinesByManufacturer", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Vaccines"); found {
			output = datasets.GroupMeasurementsByType(apiResult, measurement.GroupByManufacturer)
		} else {
			err2 = fmt.Errorf("cache does not contain Vaccines entries")
		}
		return
	})
}
