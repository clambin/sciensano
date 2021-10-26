package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"sort"
	"time"
)

// Vaccinations is the response of the GetVaccinations
type Vaccinations struct {
	Timestamps []time.Time
	Groups     []GroupedVaccinationsEntry
}

// GroupedVaccinationsEntry contains the values for the (grouped) vaccinations
type GroupedVaccinationsEntry struct {
	Name   string
	Values []*VaccinationsEntry
}

// VaccinationsEntry contains the vaccination values for a single timestamp
type VaccinationsEntry struct {
	Partial    int
	Full       int
	SingleDose int
	Booster    int
}

// Total calculates the total number of vaccinations for one entry
func (entry VaccinationsEntry) Total() int {
	return entry.Partial + entry.Full + entry.SingleDose + entry.Booster
}

// VaccinationGetter contains all required methods to retrieve vaccination data
type VaccinationGetter interface {
	GetVaccinations(ctx context.Context, end time.Time) (results *Vaccinations, err error)
	GetVaccinationsByAgeGroup(ctx context.Context, end time.Time) (results *Vaccinations, err error)
	GetVaccinationsByRegion(ctx context.Context, endTime time.Time) (results *Vaccinations, err error)
}

const (
	groupVaccinationsByNone = iota
	groupVaccinationsByAgeGroup
	groupVaccinationsByRegion
)

// GetVaccinations returns all vaccinations up to endTime
func (client *Client) GetVaccinations(ctx context.Context, endTime time.Time) (results *Vaccinations, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse
	if apiResult, err = client.Getter.GetVaccinations(ctx); err != nil {
		return
	}
	return groupVaccinations(apiResult, endTime, groupVaccinationsByNone), nil
}

// GetVaccinationsByAgeGroup returns all vaccinations, grouped by age group, up to endTime.
func (client *Client) GetVaccinationsByAgeGroup(ctx context.Context, endTime time.Time) (results *Vaccinations, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse
	if apiResult, err = client.Getter.GetVaccinations(ctx); err != nil {
		return
	}
	return groupVaccinations(apiResult, endTime, groupVaccinationsByAgeGroup), nil
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region, up to endTime.
func (client *Client) GetVaccinationsByRegion(ctx context.Context, endTime time.Time) (results *Vaccinations, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse
	if apiResult, err = client.Getter.GetVaccinations(ctx); err != nil {
		return
	}
	return groupVaccinations(apiResult, endTime, groupVaccinationsByRegion), nil
}

func groupVaccinations(vaccinations []*apiclient.APIVaccinationsResponse, endTime time.Time, groupField int) (results *Vaccinations) {
	mappedVaccinations, mappedGroups := mapVaccinations(vaccinations, endTime, groupField)
	timestamps := getUniqueSortedTimestampsFromVaccinations(mappedVaccinations)
	groups := getUniqueSortedGroupNames(mappedGroups)

	results = &Vaccinations{
		Timestamps: timestamps,
		Groups:     make([]GroupedVaccinationsEntry, len(groups)),
	}

	for index, group := range groups {
		results.Groups[index] = GroupedVaccinationsEntry{
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

func mapVaccinations(vaccinations []*apiclient.APIVaccinationsResponse, endTime time.Time, groupField int) (mappedVaccinations map[time.Time]map[string]*VaccinationsEntry, mappedGroups map[string]struct{}) {
	mappedVaccinations = make(map[time.Time]map[string]*VaccinationsEntry)
	mappedGroups = make(map[string]struct{})

	for _, entry := range vaccinations {
		if entry.TimeStamp.After(endTime) {
			continue
		}

		mappedVaccinationEntries, ok := mappedVaccinations[entry.TimeStamp.Time]
		if ok == false {
			mappedVaccinationEntries = make(map[string]*VaccinationsEntry)
		}

		var groupName string
		switch groupField {
		case groupVaccinationsByNone:
			groupName = ""
		case groupVaccinationsByAgeGroup:
			groupName = entry.AgeGroup
		case groupVaccinationsByRegion:
			groupName = entry.Region
		}

		value, ok := mappedVaccinationEntries[groupName]
		if ok == false {
			value = &VaccinationsEntry{}
		}
		switch entry.Dose {
		case "A":
			value.Partial += entry.Count
		case "B":
			value.Full += entry.Count
		case "C":
			value.SingleDose += entry.Count
		case "E":
			value.Booster += entry.Count
		}
		mappedVaccinationEntries[groupName] = value
		mappedVaccinations[entry.TimeStamp.Time] = mappedVaccinationEntries

		mappedGroups[groupName] = struct{}{}
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
func AccumulateVaccinations(vaccinationData *Vaccinations) {
	for _, group := range vaccinationData.Groups {
		partial := 0
		full := 0
		singleDose := 0
		booster := 0

		for _, value := range group.Values {
			partial += value.Partial
			value.Partial = partial

			full += value.Full
			value.Full = full

			singleDose += value.SingleDose
			value.SingleDose = singleDose

			booster += value.Booster
			value.Booster = booster
		}
	}
}
