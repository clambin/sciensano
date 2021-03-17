package mockapi

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"time"
)

var DefaultTests = []sciensano.Test{
	{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Total: 0, Positive: 0},
	{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Total: 1, Positive: 0},
	{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Total: 2, Positive: 1},
	{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Total: 3, Positive: 3},
	{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), Total: 9, Positive: 4},
	{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), Total: 10, Positive: 5},
	{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), Total: 15, Positive: 7},
}

var DefaultVaccinations = []sciensano.Vaccination{
	{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), FirstDose: 0, SecondDose: 0},
	{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), FirstDose: 1, SecondDose: 0},
	{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), FirstDose: 2, SecondDose: 1},
	{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), FirstDose: 3, SecondDose: 2},
	{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), FirstDose: 4, SecondDose: 3},
	{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), FirstDose: 5, SecondDose: 4},
	{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), FirstDose: 6, SecondDose: 5},
}

type API struct {
	Tests        []sciensano.Test
	Vaccinations []sciensano.Vaccination
}

func (client *API) GetTests(end time.Time) (results []sciensano.Test, err error) {
	for _, test := range client.Tests {
		if test.Timestamp.After(end) == false {
			results = append(results, test)
		}
	}
	return
}

func (client *API) GetVaccinations(end time.Time) (results []sciensano.Vaccination, err error) {
	for _, vaccination := range client.Vaccinations {
		if vaccination.Timestamp.After(end) == false {
			results = append(results, vaccination)
		}
	}
	return
}

func (client *API) GetVaccinationsByAge(end time.Time) (results map[string][]sciensano.Vaccination, err error) {
	results = make(map[string][]sciensano.Vaccination)
	results["45-54"], err = client.GetVaccinations(end)
	return
}

func (client *API) GetVaccinationsByRegion(end time.Time) (results map[string][]sciensano.Vaccination, err error) {
	results = make(map[string][]sciensano.Vaccination)
	results["Flanders"], err = client.GetVaccinations(end)
	return
}
