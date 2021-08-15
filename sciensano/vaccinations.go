package sciensano

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

type Vaccination struct {
	Timestamp  time.Time
	FirstDose  int
	SecondDose int
}

func (client *Client) GetVaccinations(ctx context.Context, end time.Time) (results []Vaccination, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.vaccinationsCache.GetVaccinations(ctx); err == nil {
		results = groupVaccinations(apiResult, end)
	}

	return
}

func (client *Client) GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.vaccinationsCache.GetVaccinations(ctx); err == nil {
		results = make(map[string][]Vaccination)
		ageGroups := getAgeGroups(apiResult)
		responses := make(map[string]chan []Vaccination)

		for _, ageGroup := range ageGroups {
			responses[ageGroup] = make(chan []Vaccination)
			go func(ageGroup string, response chan []Vaccination) {
				result := filterByAgeGroup(apiResult, ageGroup)
				response <- groupVaccinations(result, end)
			}(ageGroup, responses[ageGroup])
		}
		for _, ageGroup := range ageGroups {
			results[ageGroup] = <-responses[ageGroup]
		}
	}

	return
}

func (client *Client) GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.vaccinationsCache.GetVaccinations(ctx); err == nil {
		results = make(map[string][]Vaccination)
		regions := getRegions(apiResult)
		responses := make(map[string]chan []Vaccination)

		for _, region := range regions {
			responses[region] = make(chan []Vaccination)
			go func(region string, response chan []Vaccination) {
				result := filterByRegion(apiResult, region)
				response <- groupVaccinations(result, end)
			}(region, responses[region])
		}
		for _, region := range regions {
			results[region] = <-responses[region]
		}
	}

	return
}

type apiVaccinationsResponse struct {
	TimeStamp string `json:"DATE"`
	Region    string `json:"REGION"`
	AgeGroup  string `json:"AGEGROUP"`
	Gender    string `json:"SEX"`
	Dose      string `json:"DOSE"`
	Count     int    `json:"Count"`
}

func filterByAgeGroup(apiResult []apiVaccinationsResponse, ageGroup string) (output []apiVaccinationsResponse) {
	for _, result := range apiResult {
		if result.AgeGroup == ageGroup {
			output = append(output, result)
		}
	}
	return
}

func filterByRegion(apiResult []apiVaccinationsResponse, region string) (output []apiVaccinationsResponse) {
	for _, result := range apiResult {
		if result.Region == region {
			output = append(output, result)
		}
	}
	return
}

func groupVaccinations(apiResult []apiVaccinationsResponse, end time.Time) (totals []Vaccination) {
	// Store the totals in a map
	accumTotal := make(map[time.Time]Vaccination)
	for _, entry := range apiResult {
		if ts, err2 := time.Parse("2006-01-02", entry.TimeStamp); err2 == nil {
			// Skip anything after the specified end date
			if ts.After(end) {
				continue
			}

			// create a running total for each timestamp using the accumTotal map
			var current Vaccination
			var ok bool
			if current, ok = accumTotal[ts]; ok == false {
				current.Timestamp = ts
			}
			switch entry.Dose {
			case "A":
				current.FirstDose += entry.Count
			case "B":
				current.SecondDose += entry.Count
			}
			accumTotal[ts] = current
		} else {
			log.WithFields(log.Fields{"err": err2, "timestamp": entry.TimeStamp}).Warning("could not parse timestamp from API. skipping entry")
		}
	}
	// For each entry in the map, create an entry in the results slice
	for _, entry := range accumTotal {
		totals = append(totals, entry)
	}
	// Maps are iterated in random order. Sort the final slice
	sort.Slice(totals, func(i, j int) bool { return totals[i].Timestamp.Before(totals[j].Timestamp) })

	return
}

func AccumulateVaccinations(entries []Vaccination) (totals []Vaccination) {
	first := 0
	second := 0

	totals = make([]Vaccination, len(entries))
	for index, entry := range entries {
		first += entry.FirstDose
		entry.FirstDose = first
		second += entry.SecondDose
		entry.SecondDose = second
		totals[index] = entry
	}
	return
}

func getAgeGroups(results []apiVaccinationsResponse) (ageGroups []string) {
	groups := make(map[string]bool)

	for _, result := range results {
		groups[result.AgeGroup] = true
	}

	for group := range groups {
		ageGroups = append(ageGroups, group)
	}
	return
}

func getRegions(results []apiVaccinationsResponse) (regions []string) {
	groups := make(map[string]bool)

	for _, result := range results {
		groups[result.Region] = true
	}

	for group := range groups {
		regions = append(regions, group)
	}
	return
}
