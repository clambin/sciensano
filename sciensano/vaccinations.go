package sciensano

import (
	"context"
	"fmt"
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
	var apiResult []*apiVaccinationsResponse

	if apiResult, err = client.vaccinationsCache.GetVaccinations(ctx); err == nil {
		results = groupVaccinations(apiResult, end)
	}

	return
}

func (client *Client) GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error) {
	var apiResult map[string][]*apiVaccinationsResponse

	apiResult, err = client.vaccinationsCache.GetVaccinationsByAge(ctx)

	if err != nil {
		return
	}

	results = make(map[string][]Vaccination)
	ageGroups := getAgeGroups(apiResult)
	responses := make(map[string]chan []Vaccination)

	for _, ageGroup := range ageGroups {
		responses[ageGroup] = make(chan []Vaccination)
		go func(ageGroup string, response chan []Vaccination) {
			response <- groupVaccinations(apiResult[ageGroup], end)
		}(ageGroup, responses[ageGroup])
	}
	for _, ageGroup := range ageGroups {
		results[ageGroup] = <-responses[ageGroup]
	}

	return
}

func (client *Client) GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error) {
	var apiResult map[string][]*apiVaccinationsResponse

	apiResult, err = client.vaccinationsCache.GetVaccinationsByRegion(ctx)

	if err != nil {
		return
	}

	results = make(map[string][]Vaccination)
	regions := getRegions(apiResult)
	responses := make(map[string]chan []Vaccination)

	for _, region := range regions {
		responses[region] = make(chan []Vaccination)
		go func(region string, response chan []Vaccination) {
			response <- groupVaccinations(apiResult[region], end)
		}(region, responses[region])
	}
	for _, region := range regions {
		results[region] = <-responses[region]
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

func parseDate(date string) (output time.Time, err error) {
	year := (((int(date[0])-'0')*10+int(date[1])-'0')*10+int(date[2])-'0')*10 + int(date[3]) - '0'
	month := time.Month((int(date[5])-'0')*10 + int(date[6]) - '0')
	day := (int(date[8])-'0')*10 + int(date[9]) - '0'

	if year == 0 || month == 0 || day == 0 {
		err = fmt.Errorf("invalid timestamp: %s", date)
	} else {
		output = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}

	return
}

func groupVaccinations(apiResult []*apiVaccinationsResponse, end time.Time) (totals []Vaccination) {
	// Store the totals in a map
	accumTotal := make(map[time.Time]Vaccination)
	for _, entry := range apiResult {
		// time.Parse needs to be called a *lot* and it's fairly slow. replace it with a more specific, faster version.
		// ts, err2 := time.Parse("2006-01-02", entry.TimeStamp)
		ts, err2 := parseDate(entry.TimeStamp)

		if err2 != nil {
			log.WithFields(log.Fields{"err": err2, "timestamp": entry.TimeStamp}).Warning("could not parse timestamp from API. skipping entry")
			continue
		}

		// Skip anything after the specified end date
		if ts.After(end) {
			continue
		}

		// create a running total for each timestamp using the accumTotal map
		current, ok := accumTotal[ts]
		if ok == false {
			current.Timestamp = ts
		}

		switch entry.Dose {
		case "A":
			current.FirstDose += entry.Count
		case "B":
			current.SecondDose += entry.Count
		}
		accumTotal[ts] = current
	}

	// For each entry in the map, create an entry in the results slice
	totals = make([]Vaccination, 0, len(accumTotal))
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

func getAgeGroups(results map[string][]*apiVaccinationsResponse) (ageGroups []string) {
	groups := make(map[string]struct{})

	for result := range results {
		groups[result] = struct{}{}
	}

	for group := range groups {
		ageGroups = append(ageGroups, group)
	}
	return
}

func getRegions(results map[string][]*apiVaccinationsResponse) (regions []string) {
	groups := make(map[string]struct{})

	for result := range results {
		groups[result] = struct{}{}
	}

	for group := range groups {
		regions = append(regions, group)
	}
	return
}
