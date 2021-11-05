package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// VaccinationGetter contains all required methods to retrieve vaccination data
type VaccinationGetter interface {
	GetVaccinations(ctx context.Context) (results *datasets.Dataset, err error)
	GetVaccinationsByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
	GetVaccinationsByRegion(ctx context.Context) (results *datasets.Dataset, err error)
}

const (
	groupVaccinationsByNone = iota
	groupVaccinationsByAgeGroup
	groupVaccinationsByRegion
)

// GetVaccinations returns all vaccinations
func (client *Client) GetVaccinations(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getVaccinations(ctx, "GetVaccinations", "Vaccinations", groupVaccinationsByNone)
}

// GetVaccinationsByAgeGroup returns all vaccinations, grouped by age group
func (client *Client) GetVaccinationsByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getVaccinations(ctx, "GetVaccinationsByAgeGroup", "VaccinationsByAgeGroup", groupVaccinationsByAgeGroup)
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region
func (client *Client) GetVaccinationsByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getVaccinations(ctx, "GetVaccinationByRegion", "VaccinationsByRegion", groupVaccinationsByRegion)
}

func (client *Client) getVaccinations(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult apiclient.APIVaccinationsResponse
		if apiResult, err = client.Getter.GetVaccinations(ctx); err == nil {
			entry.Data = groupVaccinations(apiResult, mode)
			client.cache.Save(cacheEntryName, entry)
		} else {
			client.cache.Clear(cacheEntryName)
		}
	})
	if err == nil && entry.Data != nil {
		results = entry.Data.Copy()
	}
	return
}

func groupVaccinations(vaccinations apiclient.APIVaccinationsResponse, groupField int) (results *datasets.Dataset) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("groupVaccinations")

	mappedVaccinations, mappedGroups := mapVaccinations(vaccinations, groupField)
	timestamps := getUniqueSortedTimestampsFromVaccinations(mappedVaccinations)
	groups := getUniqueSortedGroupNames(mappedGroups)

	results = &datasets.Dataset{
		Timestamps: timestamps,
		Groups:     make([]datasets.GroupedDatasetEntry, len(groups)),
	}

	for index, group := range groups {
		results.Groups[index] = datasets.GroupedDatasetEntry{
			Name: group,
		}
	}

	for _, timestamp := range timestamps {
		for index, group := range groups {
			entry, ok := mappedVaccinations[timestamp][group]
			if ok == false {
				entry = &VaccinationsEntry{}
			}
			results.Groups[index].Values = append(results.Groups[index].Values, entry)
		}
	}
	return
}

func getUniqueSortedTimestampsFromVaccinations(input map[time.Time]map[string]*VaccinationsEntry) (output []time.Time) {
	for timestamp := range input {
		output = append(output, timestamp)
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Before(output[j]) })
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

func mapVaccinations(vaccinations apiclient.APIVaccinationsResponse, groupField int) (mappedVaccinations map[time.Time]map[string]*VaccinationsEntry, mappedGroups map[string]struct{}) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("mapVaccinations")

	if vaccinations == nil || len(vaccinations) == 0 {
		return
	}
	// sort.Slice(vaccinations, func(i, j int) bool { return vaccinations[i].TimeStamp.Time.Before(vaccinations[j].TimeStamp.Time)} )

	mappedVaccinations = make(map[time.Time]map[string]*VaccinationsEntry)
	mappedGroups = make(map[string]struct{})

	currentTimestamp := time.Time{}
	currentEntries := make(map[string]*VaccinationsEntry)

	for _, vaccination := range vaccinations {
		if !currentTimestamp.IsZero() && !currentTimestamp.Equal(vaccination.TimeStamp.Time) {
			mappedVaccinations[currentTimestamp] = currentEntries
			currentEntries = make(map[string]*VaccinationsEntry)
		}
		currentTimestamp = vaccination.TimeStamp.Time

		var groupName string
		switch groupField {
		case groupVaccinationsByNone:
			groupName = ""
		case groupVaccinationsByAgeGroup:
			groupName = vaccination.AgeGroup
		case groupVaccinationsByRegion:
			groupName = vaccination.Region
		}

		mappedGroups[groupName] = struct{}{}

		entry, ok := currentEntries[groupName]

		if ok == false {
			entry = &VaccinationsEntry{}
		}

		entry.Add(vaccination)
		currentEntries[groupName] = entry
	}

	if !currentTimestamp.IsZero() {
		mappedVaccinations[currentTimestamp] = currentEntries
	}

	return
}
