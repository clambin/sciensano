package sciensano

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"time"
)

type Vaccination struct {
	Timestamp  time.Time
	FirstDose  int
	SecondDose int
}

type Vaccinations []Vaccination

func GetVaccinationsTargets() []string {
	return []string{
		"vaccinations-first",
		"vaccinations-second",
	}
}

var AgeGroups = []string{
	"",
	"0-17",
	"18-34",
	"35-44",
	"45-54",
	"55-64",
	"65-74",
	"75-84",
	"85+",
}

func GetVaccinationsByAgeTargets() (targets []string) {
	for _, ageGroup := range AgeGroups {
		targets = append(targets, "vaccinations-"+ageGroup+"-first")
		targets = append(targets, "vaccinations-"+ageGroup+"-second")
	}
	return
}

func (client *Client) GetVaccinations(end time.Time) (results Vaccinations, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.getVaccinations(); err == nil {
		// debug: are we missing something
		sizes := make(map[string][]apiVaccinationsResponse)
		for _, ageGroup := range AgeGroups {
			sizes[ageGroup] = filterByAgeGroup(apiResult, ageGroup)
		}

		results = accumulateVaccinations(groupVaccinations(apiResult, end))
	}

	return
}

func (client *Client) GetVaccinationsByAge(end time.Time, group string) (results Vaccinations, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.getVaccinations(); err == nil {
		apiResult = filterByAgeGroup(apiResult, group)
		results = accumulateVaccinations(groupVaccinations(apiResult, end))
	}

	return
}

func GetAgeGroupFromTarget(target string) (output string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		output = strings.TrimPrefix(target, "vaccinations-")
		output = strings.TrimSuffix(output, "-first")
		output = strings.TrimSuffix(output, "-second")
	}
	return
}

func GetModeFromTarget(target string) (mode string) {
	if strings.HasSuffix(target, "-first") {
		mode = "A"
	} else if strings.HasSuffix(target, "-second") {
		mode = "B"
	}
	return
}

func filterByAgeGroup(apiResult []apiVaccinationsResponse, ageGroup string) (output []apiVaccinationsResponse) {
	for _, result := range apiResult {
		if result.AgeGroup == ageGroup {
			output = append(output, result)
		}
	}
	return
}

func groupVaccinations(apiResult []apiVaccinationsResponse, end time.Time) (totals Vaccinations) {
	// Store the totals in a map
	accumTotal := make(map[time.Time]Vaccination, 0)
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
				current = Vaccination{Timestamp: ts}
			}
			if entry.Dose == "A" {
				current.FirstDose += entry.Count
			} else if entry.Dose == "B" {
				current.SecondDose += entry.Count
			} else {
				log.WithField("dose", entry.Dose).Warning("unexpected dose code. skipping entry")
				continue
			}
			accumTotal[ts] = current
		} else {
			log.WithFields(log.Fields{
				"err":       err2,
				"timestamp": entry.TimeStamp,
			}).Warning("could not parse timestamp from API. skipping entry")
		}
	}
	// For each entry in the map, create an entry in the results slice
	for _, entry := range accumTotal {
		totals = append(totals, entry)
	}
	// Maps are iterated in random order. Sort the final slice
	sort.Sort(totals)

	// log.WithField("ageGroups", ageGroups).Info("ageGroups")

	return
}

// helper functions for sort.Sort([]Test)
func (p Vaccinations) Len() int {
	return len(p)
}

func (p Vaccinations) Less(i, j int) bool {
	return p[i].Timestamp.Before(p[j].Timestamp)
}

func (p Vaccinations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func accumulateVaccinations(entries Vaccinations) (totals Vaccinations) {
	first := 0
	second := 0

	totals = make(Vaccinations, len(entries))
	for index, entry := range entries {
		first += entry.FirstDose
		second += entry.SecondDose
		totals[index] = Vaccination{
			Timestamp:  entry.Timestamp,
			FirstDose:  first,
			SecondDose: second,
		}
	}
	return
}
