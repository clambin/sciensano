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
	// Partial is how many people were partially vaccinated on this day
	Partial int
	// Full is how many people were fully vaccinated on this day
	Full int
}

// GetVaccinations returns all vaccinations up to endTime
func (client *Client) GetVaccinations(ctx context.Context, endTime time.Time) (results []Vaccination, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse

	if apiResult, err = client.APIClient.GetVaccinations(ctx); err == nil {
		results = groupVaccinations(apiResult, endTime, true)
	}

	return
}

// GetVaccinationsForLag returns all vaccinations up to endTime so the caller can determine the lag between partial and
// full vaccination, i.e. exclude any vaccines that only need one dose
func (client *Client) GetVaccinationsForLag(ctx context.Context, endTime time.Time) (results []Vaccination, err error) {
	var apiResult []*apiclient.APIVaccinationsResponse

	if apiResult, err = client.APIClient.GetVaccinations(ctx); err == nil {
		results = groupVaccinations(apiResult, endTime, false)
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
		results[ageGroup] = groupVaccinations(grouped[ageGroup], end, true)
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
		results[region] = groupVaccinations(grouped[region], end, true)
	}

	return
}

func groupVaccinations(apiResult []*apiclient.APIVaccinationsResponse, end time.Time, includeSingleDoseVaccines bool) (totals []Vaccination) {
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
			if includeSingleDoseVaccines {
				current.Full += entry.Count
			}
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
		first += entry.Partial
		entry.Partial = first
		second += entry.Full
		entry.Full = second
		totals[index] = entry
	}
	return
}
