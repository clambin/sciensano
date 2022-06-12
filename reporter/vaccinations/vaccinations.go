package vaccinations

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	apiCache "github.com/clambin/sciensano/apiclient/cache"
	"github.com/clambin/sciensano/apiclient/sciensano"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
	"strconv"
)

const (
	// TypeAll tells getByType not to filter by type of vaccination
	TypeAll = iota
	// TypePartial filters partial vaccinations
	TypePartial
	// TypeFull filters full vaccinations. It counts 2nd vaccinations and single dose vaccinations
	TypeFull
	// TypeBooster filters booster vaccinations
	TypeBooster
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APICache    apiCache.Holder
}

// Get returns all vaccinations
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("Vaccinations", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("Vaccinations"); found {
			output = table.NewFromAPIResponse(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Vaccinations entries")
		}
		return output, err2
	})
}

// GetByAgeGroup returns all vaccinations, grouped by age group
func (r *Reporter) GetByAgeGroup(vaccinationType int) (results *data.Table, err error) {
	name := "VaccinationsByAgeGroup-" + strconv.Itoa(vaccinationType)
	return r.ReportCache.MaybeGenerate(name, func() (*data.Table, error) {
		return r.getByType(apiclient.GroupByAgeGroup, vaccinationType)
	})
}

// GetByRegion returns all vaccinations, grouped by region
func (r *Reporter) GetByRegion(vaccinationType int) (results *data.Table, err error) {
	name := "VaccinationsByRegion-" + strconv.Itoa(vaccinationType)
	return r.ReportCache.MaybeGenerate(name, func() (*data.Table, error) {
		return r.getByType(apiclient.GroupByRegion, vaccinationType)
	})
}

// GetByManufacturer returns all vaccinations, grouped by manufacturer
func (r *Reporter) GetByManufacturer() (results *data.Table, err error) {
	name := "VaccinationsByManufacturer"
	return r.ReportCache.MaybeGenerate(name, func() (*data.Table, error) {
		return r.getByType(apiclient.GroupByManufacturer, TypeAll)
	})
}

func (r *Reporter) getByType(mode int, vaccinationType int) (results *data.Table, err error) {
	apiResult, found := r.APICache.Get("Vaccinations")
	if !found {
		err = fmt.Errorf("cache does not contain Vaccinations entries")
	}

	if vaccinationType != TypeAll {
		apiResult = filterVaccinations(apiResult, vaccinationType)
	}
	results = table.NewGroupedFromAPIResponse(apiResult, mode)
	return
}

func filterVaccinations(input []apiclient.APIResponse, vaccinationType int) (output []apiclient.APIResponse) {
	output = make([]apiclient.APIResponse, 0, len(input))
	for _, entry := range input {
		// this is faster than using cache.GetAttributeValues()
		dose := entry.(*sciensano.APIVaccinationsResponse).Dose
		if (vaccinationType == TypePartial && dose == "A") ||
			(vaccinationType == TypeFull && (dose == "B" || dose == "C")) ||
			(vaccinationType == TypeBooster && dose == "E") {
			output = append(output, entry)
		}
	}
	return
}
