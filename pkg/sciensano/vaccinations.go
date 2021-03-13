package sciensano

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

var AgeGroups = []string{
	"0-17",
	"18-34",
	"35-44",
	"45-54",
	"55-64",
	"65-74",
	"75-84",
	"85+",
}

var Regions = []string{
	"",
	"Flanders",
	"Brussels",
	"Wallonia",
	"Ostbelgien",
}

type Vaccination struct {
	Timestamp  time.Time
	FirstDose  int
	SecondDose int
}

func (client *Client) GetVaccinations(end time.Time) (results []Vaccination, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.getVaccinations(); err == nil {
		results = groupVaccinations(apiResult, end)
	}

	return
}

func (client *Client) GetVaccinationsByAge(end time.Time, group string) (results []Vaccination, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.getVaccinations(); err == nil {
		apiResult = filterByAgeGroup(apiResult, group)
		results = groupVaccinations(apiResult, end)
	}

	return
}

func (client *Client) GetVaccinationsByRegion(end time.Time, group string) (results []Vaccination, err error) {
	var apiResult []apiVaccinationsResponse

	if apiResult, err = client.getVaccinations(); err == nil {
		apiResult = filterByRegion(apiResult, group)
		results = groupVaccinations(apiResult, end)
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

func (client *Client) getVaccinations() ([]apiVaccinationsResponse, error) {
	var err error

	if client.vaccinationsCache == nil || time.Now().After(client.vaccinationsCacheExpiry) {

		req, _ := http.NewRequest("GET", baseURL+"COVID19BE_VACC.json", nil)

		var resp *http.Response
		var stats []apiVaccinationsResponse

		if resp, err = client.apiClient.Do(req); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				var (
					body []byte
				)
				if body, err = ioutil.ReadAll(resp.Body); err == nil {
					if err = json.Unmarshal(body, &stats); err == nil {
						client.vaccinationsCache = stats
						client.vaccinationsCacheExpiry = time.Now().Add(client.VaccinationsCacheDuration)
					}
				}
			} else {
				err = errors.New(resp.Status)
			}
		}
	}

	return client.vaccinationsCache, err
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
			switch entry.Dose {
			case "A":
				current.FirstDose += entry.Count
			case "B":
				current.SecondDose += entry.Count
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

	return
}

func AccumulateVaccinations(entries []Vaccination) (totals []Vaccination) {
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

// helper functions for sort.Sort(Vaccinations)
type Vaccinations []Vaccination

func (p Vaccinations) Len() int {
	return len(p)
}

func (p Vaccinations) Less(i, j int) bool {
	return p[i].Timestamp.Before(p[j].Timestamp)
}

func (p Vaccinations) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}