package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/simplejson/v3/dataset"
	"strconv"
)

type VaccinationType int

const (
	// VaccinationTypeAll tells realGetVaccinationByType not to filter by type of vaccination
	VaccinationTypeAll = iota
	// VaccinationTypePartial filters partial vaccinations
	VaccinationTypePartial
	// VaccinationTypeFull filters full vaccinations. It counts 2nd vaccinations and single dose vaccinations
	VaccinationTypeFull
	// VaccinationTypeBooster filters booster vaccinations
	VaccinationTypeBooster
)

// VaccinationGetter contains all required methods to retrieve vaccination data
type VaccinationGetter interface {
	GetVaccinations() (results *dataset.Dataset, err error)
	GetVaccinationsByAgeGroup(vaccinationType VaccinationType) (results *dataset.Dataset, err error)
	GetVaccinationsByRegion(vaccinationType VaccinationType) (results *dataset.Dataset, err error)
	GetVaccinationsByManufacturer() (results *dataset.Dataset, err error)
}

// GetVaccinations returns all vaccinations
func (client *Client) GetVaccinations() (results *dataset.Dataset, err error) {
	return client.ReportCache.MaybeGenerate("Vaccinations", func() (output *dataset.Dataset, err2 error) {
		if apiResult, found := client.APICache.Get("Vaccinations"); found {
			output = NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Vaccinations entries")
		}
		return output, err2
	})
}

// GetVaccinationsByAgeGroup returns all vaccinations, grouped by age group
func (client *Client) GetVaccinationsByAgeGroup(vaccinationType VaccinationType) (results *dataset.Dataset, err error) {
	name := "VaccinationsByAgeGroup-" + strconv.Itoa(int(vaccinationType))
	return client.ReportCache.MaybeGenerate(name, func() (*dataset.Dataset, error) {
		return client.realGetVaccinationByType(apiclient.GroupByAgeGroup, vaccinationType)
	})
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region
func (client *Client) GetVaccinationsByRegion(vaccinationType VaccinationType) (results *dataset.Dataset, err error) {
	name := "VaccinationsByRegion-" + strconv.Itoa(int(vaccinationType))
	return client.ReportCache.MaybeGenerate(name, func() (*dataset.Dataset, error) {
		return client.realGetVaccinationByType(apiclient.GroupByRegion, vaccinationType)
	})
}

// GetVaccinationsByManufacturer returns all vaccinations, grouped by manufacturer
func (client *Client) GetVaccinationsByManufacturer() (results *dataset.Dataset, err error) {
	name := "VaccinationsByManufacturer"
	return client.ReportCache.MaybeGenerate(name, func() (*dataset.Dataset, error) {
		return client.realGetVaccinationByType(apiclient.GroupByManufacturer, VaccinationTypeAll)
	})
}

func (client *Client) realGetVaccinationByType(mode int, vaccinationType VaccinationType) (results *dataset.Dataset, err error) {
	if apiResult, found := client.APICache.Get("Vaccinations"); found {
		if vaccinationType != VaccinationTypeAll {
			apiResult = filterVaccinations(apiResult, vaccinationType)
		}
		results = NewGroupedFromAPIResponse(apiResult, mode)
	} else {
		err = fmt.Errorf("cache does not contain Vaccinations entries")
	}
	return
}

func filterVaccinations(input []apiclient.APIResponse, vaccinationType VaccinationType) (output []apiclient.APIResponse) {
	output = make([]apiclient.APIResponse, 0, len(input))
	for _, entry := range input {
		// this is faster than using cache.GetAttributeValues()
		dose := entry.(*sciensano.APIVaccinationsResponse).Dose
		if (vaccinationType == VaccinationTypePartial && dose == "A") ||
			(vaccinationType == VaccinationTypeFull && (dose == "B" || dose == "C")) ||
			(vaccinationType == VaccinationTypeBooster && dose == "E") {
			output = append(output, entry)
		}
	}
	return
}
