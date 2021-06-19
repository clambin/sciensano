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
	{Timestamp: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC), Total: 20, Positive: 8},
	{Timestamp: time.Date(2021, 1, 9, 0, 0, 0, 0, time.UTC), Total: 21, Positive: 8},
	{Timestamp: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC), Total: 22, Positive: 10},
	{Timestamp: time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC), Total: 23, Positive: 12},
	{Timestamp: time.Date(2021, 1, 12, 0, 0, 0, 0, time.UTC), Total: 24, Positive: 14},
	{Timestamp: time.Date(2021, 1, 13, 0, 0, 0, 0, time.UTC), Total: 25, Positive: 15},
	{Timestamp: time.Date(2021, 1, 14, 0, 0, 0, 0, time.UTC), Total: 26, Positive: 18},
	{Timestamp: time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC), Total: 27, Positive: 14},
	{Timestamp: time.Date(2021, 1, 16, 0, 0, 0, 0, time.UTC), Total: 28, Positive: 12},
}

var DefaultVaccinations = []sciensano.Vaccination{
	{Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), FirstDose: 0, SecondDose: 0},
	{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), FirstDose: 1, SecondDose: 0},
	{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), FirstDose: 2, SecondDose: 1},
	{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), FirstDose: 3, SecondDose: 2},
	{Timestamp: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), FirstDose: 4, SecondDose: 3},
	{Timestamp: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), FirstDose: 5, SecondDose: 4},
	{Timestamp: time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), FirstDose: 6, SecondDose: 5},
	{Timestamp: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC), FirstDose: 7, SecondDose: 6},
	{Timestamp: time.Date(2021, 1, 9, 0, 0, 0, 0, time.UTC), FirstDose: 8, SecondDose: 7},
	{Timestamp: time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC), FirstDose: 9, SecondDose: 8},
	{Timestamp: time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC), FirstDose: 10, SecondDose: 9},
	{Timestamp: time.Date(2021, 1, 12, 0, 0, 0, 0, time.UTC), FirstDose: 11, SecondDose: 10},
	{Timestamp: time.Date(2021, 1, 13, 0, 0, 0, 0, time.UTC), FirstDose: 12, SecondDose: 11},
	{Timestamp: time.Date(2021, 1, 14, 0, 0, 0, 0, time.UTC), FirstDose: 13, SecondDose: 12},
	{Timestamp: time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC), FirstDose: 14, SecondDose: 13},
	{Timestamp: time.Date(2021, 1, 16, 0, 0, 0, 0, time.UTC), FirstDose: 15, SecondDose: 14},
}

var AltVaccinations = []sciensano.Vaccination{
	{Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), FirstDose: 250, SecondDose: 25},
	{Timestamp: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), FirstDose: 275, SecondDose: 50},
	{Timestamp: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), FirstDose: 375, SecondDose: 150},
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
