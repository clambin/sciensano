package reporter

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/simplejson/v3/dataset"
)

// VaccinesGetter contains all required methods to retrieve vaccines data
type VaccinesGetter interface {
	GetVaccines() (results *dataset.Dataset, err error)
	GetVaccinesByManufacturer() (results *dataset.Dataset, err error)
}

// GetVaccines returns all vaccines data
func (client *Client) GetVaccines() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("Vaccines", func() (*dataset.Dataset, error) {
		apiResult, ok := client.APICache.Get("Vaccines")

		if ok == false {
			return nil, errors.New("cache does not contain Vaccines entries")
		}
		return NewFromAPIResponse(apiResult), nil

	})
}

// GetVaccinesByManufacturer returns all hospitalisations, grouped by manufacturer
func (client *Client) GetVaccinesByManufacturer() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("VaccinesByManufacturer", func() (*dataset.Dataset, error) {
		apiResult, ok := client.APICache.Get("Vaccines")

		if ok == false {
			return nil, errors.New("cache does not contain Vaccines entries")
		}
		return NewGroupedFromAPIResponse(apiResult, apiclient.GroupByManufacturer), nil
	})
}
