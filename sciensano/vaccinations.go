package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// VaccinationGetter contains all required methods to retrieve vaccination data
type VaccinationGetter interface {
	GetVaccinations(ctx context.Context) (results *datasets.Dataset, err error)
	GetVaccinationsByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
	GetVaccinationsByRegion(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetVaccinations returns all vaccinations
func (client *Client) GetVaccinations(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getVaccinations(ctx, "GetVaccinations", "Vaccinations", apiclient.GroupByNone)
}

// GetVaccinationsByAgeGroup returns all vaccinations, grouped by age group
func (client *Client) GetVaccinationsByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getVaccinations(ctx, "GetVaccinationsByAgeGroup", "VaccinationsByAgeGroup", apiclient.GroupByAgeGroup)
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region
func (client *Client) GetVaccinationsByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getVaccinations(ctx, "GetVaccinationByRegion", "VaccinationsByRegion", apiclient.GroupByRegion)
}

func (client *Client) getVaccinations(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.Cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult []apiclient.Measurement
		if apiResult, err = client.Getter.GetVaccinations(ctx); err == nil {
			entry.Data = groupMeasurements(apiResult, mode, NewVaccinationsEntry)
			client.Cache.Save(cacheEntryName, entry)
		} else {
			client.Cache.Clear(cacheEntryName)
		}
	})
	if err == nil && entry.Data != nil {
		results = entry.Data.Copy()
	}
	return
}

// AccumulateVaccinations takes a list of vaccinations and accumulates the doses
func AccumulateVaccinations(vaccinationData *datasets.Dataset) {
	for _, group := range vaccinationData.Groups {
		partial := 0
		full := 0
		singleDose := 0
		booster := 0

		for _, value := range group.Values {
			partial += value.(*VaccinationsEntry).Partial
			value.(*VaccinationsEntry).Partial = partial

			full += value.(*VaccinationsEntry).Full
			value.(*VaccinationsEntry).Full = full

			singleDose += value.(*VaccinationsEntry).SingleDose
			value.(*VaccinationsEntry).SingleDose = singleDose

			booster += value.(*VaccinationsEntry).Booster
			value.(*VaccinationsEntry).Booster = booster
		}
	}
}
