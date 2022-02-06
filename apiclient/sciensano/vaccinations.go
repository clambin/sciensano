package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/mailru/easyjson"
	"time"
)

// APIVaccinationsResponse is a single entry in APIVaccinationResponse
type APIVaccinationsResponse struct {
	TimeStamp    TimeStamp `json:"DATE"`
	Manufacturer string    `json:"BRAND"`
	Region       string    `json:"REGION"`
	AgeGroup     string    `json:"AGEGROUP"`
	Gender       string    `json:"SEX"`
	Dose         string    `json:"DOSE"`
	Count        int       `json:"COUNT"`
}

// APIVaccinationsResponses is a slice of APIVaccinationResponse structures
//easyjson:json
type APIVaccinationsResponses []*APIVaccinationsResponse

var _ apiclient.APIResponse = &APIVaccinationsResponse{}

// GetTimestamp returns the entry's timestamp
func (v *APIVaccinationsResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APIVaccinationsResponse) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case apiclient.GroupByRegion:
		value = v.Region
	case apiclient.GroupByAgeGroup:
		value = v.AgeGroup
	case apiclient.GroupByManufacturer:
		value = v.Manufacturer
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIVaccinationsResponse) GetTotalValue() float64 {
	return float64(v.Count)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIVaccinationsResponse) GetAttributeNames() []string {
	return []string{"partial", "full", "singledose", "booster"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIVaccinationsResponse) GetAttributeValues() (values []float64) {
	values = []float64{0, 0, 0, 0}
	switch v.Dose {
	case "A":
		values[0] = float64(v.Count)
	case "B":
		values[1] = float64(v.Count)
	case "C":
		values[2] = float64(v.Count)
	case "E":
		values[3] = float64(v.Count)
	}
	return
}

// GetVaccinations retrieves all COVID-19 vaccinations.
func (client *Client) GetVaccinations(ctx context.Context) (results []apiclient.APIResponse, err error) {
	var body []byte
	body, err = client.call(ctx, "vaccinations")
	if err != nil {
		return
	}

	var response APIVaccinationsResponses
	err = easyjson.Unmarshal(body, &response)
	if err != nil {
		return
	}

	results = make([]apiclient.APIResponse, len(response))
	for index, entry := range response {
		results[index] = entry
	}
	return
}
