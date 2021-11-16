package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
)

const (
	// VaccinationTypePartial tells GetValue to return the partial vaccination count
	VaccinationTypePartial int = iota
	// VaccinationTypeFull tells GetValue to return the full vaccination count. It counts 2nd vaccinations and single dose vaccinations
	VaccinationTypeFull
	// VaccinationTypeBooster tells GetValue to return the booster vaccination count
	VaccinationTypeBooster
)

// VaccinationGetter contains all required methods to retrieve vaccination data
type VaccinationGetter interface {
	GetVaccinations() (results *datasets.Dataset, err error)
	GetVaccinationsByAgeGroup(vaccinationType int) (results *datasets.Dataset, err error)
	GetVaccinationsByRegion(vaccinationType int) (results *datasets.Dataset, err error)
}

// GetVaccinations returns all vaccinations
func (client *Client) GetVaccinations() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Vaccinations", func() (output *datasets.Dataset, err2 error) {
		if apiResult, found := client.Sciensano.Get("Vaccinations"); found {
			output = datasets.GroupMeasurements(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Vaccinations entries")
		}
		return output, err2
	})
}

// GetVaccinationsByAgeGroup returns all vaccinations, grouped by age group
func (client *Client) GetVaccinationsByAgeGroup(vaccinationType int) (results *datasets.Dataset, err error) {
	name := fmt.Sprintf("VaccinationsByAgeGroup-%d", vaccinationType)
	return client.Cache.MaybeGenerate(name, func() (*datasets.Dataset, error) {
		return client.realGetVaccinationByType(measurement.GroupByAgeGroup, vaccinationType)
	})
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region
func (client *Client) GetVaccinationsByRegion(vaccinationType int) (results *datasets.Dataset, err error) {
	name := fmt.Sprintf("VaccinationsByRegion-%d", vaccinationType)
	return client.Cache.MaybeGenerate(name, func() (*datasets.Dataset, error) {
		return client.realGetVaccinationByType(measurement.GroupByRegion, vaccinationType)
	})
}

func (client *Client) realGetVaccinationByType(mode, vaccinationType int) (results *datasets.Dataset, err error) {
	if apiResult, found := client.Sciensano.Get("Vaccinations"); found {
		apiResult = filterVaccinations(apiResult, vaccinationType)
		results = datasets.GroupMeasurementsByType(apiResult, mode)
	} else {
		err = fmt.Errorf("cache does not contain Vaccinations entries")
	}
	return
}

func filterVaccinations(input []measurement.Measurement, vaccinationType int) (output []measurement.Measurement) {
	output = make([]measurement.Measurement, 0, len(input))
	for _, entry := range input {
		// this is faster than using measurement.GetAttributeValues()
		dose := entry.(*sciensano.APIVaccinationsResponseEntry).Dose
		if (vaccinationType == VaccinationTypePartial && dose == "A") ||
			(vaccinationType == VaccinationTypeFull && (dose == "B" || dose == "C")) ||
			(vaccinationType == VaccinationTypeBooster && dose == "E") {
			output = append(output, entry)
		}
	}
	return
}
