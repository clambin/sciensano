package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// Vaccination represents the number of first/second doses administered on a given day
type Vaccination struct {
	// Timestamp is the day of the measurement
	Timestamp time.Time
	// Partial is how many people received a first vaccine requiring two vaccinations
	Partial int
	// Full is how many people received a second vaccine requiring two vaccinations
	Full int
	// SingleDose is how many people received a vaccine requiring a single dose (e.g. J&J)
	SingleDose int
	// Booster is how many peoplse received a booster vaccine
	Booster int
}

func (vaccination Vaccination) Total() int {
	return vaccination.Partial + vaccination.Full + vaccination.SingleDose + vaccination.Booster
}

// VaccinationGetter contains all required methods to retrieve vaccination data
type VaccinationGetter interface {
	GetVaccinations(ctx context.Context, endTime time.Time) (results []Vaccination, err error)
	GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
	GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error)
}

// GetVaccinations returns all vaccinations up to endTime
func (client *Client) GetVaccinations(ctx context.Context, endTime time.Time) (results []Vaccination, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse
	apiResult, err = client.APIClient.GetVaccinations(ctx)
	if err != nil {
		return
	}

	results = groupVaccinations(apiResult, endTime)
	return
}

// GetVaccinationsByAge returns all vaccinations, grouped by age group, up to endTime.
func (client *Client) GetVaccinationsByAge(ctx context.Context, endTime time.Time) (results map[string][]Vaccination, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse
	apiResult, err = client.APIClient.GetVaccinations(ctx)
	if err != nil {
		return
	}

	grouped := make(map[string][]*apiclient.APIVaccinationsResponse)
	for _, entry := range apiResult {
		grouped[entry.AgeGroup] = append(grouped[entry.AgeGroup], entry)
	}

	results = make(map[string][]Vaccination)
	for ageGroup := range grouped {
		results[ageGroup] = groupVaccinations(grouped[ageGroup], endTime)
	}

	return
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region, up to endTime.
func (client *Client) GetVaccinationsByRegion(ctx context.Context, endTime time.Time) (results map[string][]Vaccination, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse
	apiResult, err = client.APIClient.GetVaccinations(ctx)
	if err != nil {
		return
	}

	grouped := make(map[string][]*apiclient.APIVaccinationsResponse)
	for _, entry := range apiResult {
		grouped[entry.Region] = append(grouped[entry.Region], entry)
	}

	results = make(map[string][]Vaccination)
	for region := range grouped {
		results[region] = groupVaccinations(grouped[region], endTime)
	}

	return
}

func groupVaccinations(apiResult []*apiclient.APIVaccinationsResponse, end time.Time) (totals []Vaccination) {
	// Note: this algorithm assumes apiResult is sorted by date.  If that changes, add:
	// sort.Slice(apiResult, func(i, j int) bool { return apiResult[i].TimeStamp.Time.Before(apiResult[j].TimeStamp.Time) })

	var current Vaccination
	for _, entry := range apiResult {
		// Skip anything after the specified end date
		if entry.TimeStamp.Time.After(end) {
			continue
		}

		if entry.TimeStamp.Time != current.Timestamp {
			if !current.Timestamp.IsZero() {
				totals = append(totals, current)
			}

			current = Vaccination{Timestamp: entry.TimeStamp.Time}
		}

		switch entry.Dose {
		case "A":
			current.Partial += entry.Count
		case "B":
			current.Full += entry.Count
		case "C":
			current.SingleDose += entry.Count
		case "E":
			current.Booster += entry.Count
		}
	}

	if !current.Timestamp.IsZero() {
		totals = append(totals, current)
	}

	return
}

// AccumulateVaccinations takes a list of vaccinations and accumulates the doses
func AccumulateVaccinations(entries []Vaccination) (totals []Vaccination) {
	first := 0
	second := 0
	single := 0
	booster := 0
	totals = make([]Vaccination, len(entries))
	for index, entry := range entries {
		first += entry.Partial
		second += entry.Full
		single += entry.SingleDose
		booster += entry.Booster

		totals[index] = Vaccination{
			Timestamp:  entry.Timestamp,
			Partial:    first,
			Full:       second,
			SingleDose: single,
			Booster:    booster,
		}
	}
	return
}
