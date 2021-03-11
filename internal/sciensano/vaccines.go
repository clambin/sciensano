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

type Vaccine struct {
	Timestamp  time.Time
	FirstDose  int
	SecondDose int
}

type Vaccines []Vaccine

func (client *APIClient) GetVaccines(end time.Time) (results Vaccines, err error) {
	var apiResult []apiVaccineResponse

	if apiResult, err = client.getVaccines(); err == nil {
		results = accumulateVaccines(groupVaccines(apiResult, end))
	}

	return
}

type apiVaccineResponse struct {
	TimeStamp string `json:"DATE"`
	Region    string `json:"REGION"`
	AgeGroup  string `json:"AGEGROUP"`
	Gender    string `json:"SEX"`
	Dose      string `json:"DOSE"`
	Count     int    `json:"Count"`
}

func (client *APIClient) getVaccines() (stats []apiVaccineResponse, err error) {
	req, _ := http.NewRequest("GET", baseURL+"COVID19BE_VACC.json", nil)

	var resp *http.Response
	if resp, err = client.client.Do(req); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			var (
				body []byte
			)
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				err = json.Unmarshal(body, &stats)
			}
		} else {
			err = errors.New(resp.Status)
		}
	}
	return
}

func groupVaccines(apiResult []apiVaccineResponse, end time.Time) (totals Vaccines) {
	// Store the totals in a map
	accumTotal := make(map[time.Time]Vaccine, 0)
	for _, entry := range apiResult {
		if ts, err2 := time.Parse("2006-01-02", entry.TimeStamp); err2 == nil {
			// Skip anything after the specified end date
			if ts.After(end) {
				continue
			}

			var current Vaccine
			var ok bool
			if current, ok = accumTotal[ts]; ok == false {
				current = Vaccine{Timestamp: ts}
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

	return
}

// helper functions for sort.Sort([]Test)
func (p Vaccines) Len() int {
	return len(p)
}

func (p Vaccines) Less(i, j int) bool {
	return p[i].Timestamp.Before(p[j].Timestamp)
}

func (p Vaccines) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func accumulateVaccines(entries Vaccines) (totals Vaccines) {
	first := 0
	second := 0

	totals = make(Vaccines, len(entries))
	for index, entry := range entries {
		first += entry.FirstDose
		second += entry.SecondDose
		totals[index] = Vaccine{
			Timestamp:  entry.Timestamp,
			FirstDose:  first,
			SecondDose: second,
		}
	}
	return
}
