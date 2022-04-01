package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/simplejson/v3/dataset"
)

// HospitalisationsGetter contains all methods providing COVID-19-related hospitalisation figures
type HospitalisationsGetter interface {
	GetHospitalisations() (results *dataset.Dataset, err error)
	GetHospitalisationsByRegion() (results *dataset.Dataset, err error)
	GetHospitalisationsByProvince() (results *dataset.Dataset, err error)
}

// GetHospitalisations returns all hospitalisations
func (client *Client) GetHospitalisations() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("Hospitalisations", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Hospitalisations"); found {
			output = NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}

// GetHospitalisationsByRegion returns all hospitalisations, grouped by region
func (client *Client) GetHospitalisationsByRegion() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("HospitalisationsByRegion", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Hospitalisations"); found {
			output = NewGroupedFromAPIResponse(apiResult, apiclient.GroupByRegion)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}

// GetHospitalisationsByProvince returns all hospitalisations, grouped by province
func (client *Client) GetHospitalisationsByProvince() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("HospitalisationsByProvince", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Hospitalisations"); found {
			output = NewGroupedFromAPIResponse(apiResult, apiclient.GroupByProvince)
		} else {
			err2 = fmt.Errorf("cache does not contain Hospitalisations entries")
		}
		return
	})
}
