package sciensano

import (
	"context"
	"github.com/clambin/sciensano/sciensano/apiclient"
	"time"
)

// Vaccination represents the number of first/second doses administered on a given day
type Vaccination struct {
	// Timestamp is the day of the measurement
	Timestamp time.Time
	// FirstDose is how many first doses were administered by this day
	FirstDose int
	// SecondDose is how many second (final) doses were administered by that day
	SecondDose int
}

// GetVaccinations returns all vaccinations up to endTime
func (client *Client) GetVaccinations(ctx context.Context, endTime time.Time) (results []Vaccination, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse

	if apiResult, err = client.APIClient.GetVaccinations(ctx); err == nil {
		results = groupVaccinations(apiResult, endTime)
	}

	return
}

// GetVaccinationsByAge returns all vaccinations, grouped by age group, up to endTime.
func (client *Client) GetVaccinationsByAge(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error) {
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
		results[ageGroup] = groupVaccinations(grouped[ageGroup], end)
	}

	return
}

// GetVaccinationsByRegion returns all vaccinations, grouped by region, up to endTime.
func (client *Client) GetVaccinationsByRegion(ctx context.Context, end time.Time) (results map[string][]Vaccination, err error) {
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
		results[region] = groupVaccinations(grouped[region], end)
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
			current.FirstDose += entry.Count
		case "B":
			current.SecondDose += entry.Count
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
