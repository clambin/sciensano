package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/mailru/easyjson"
	"time"
)

// APICasesResponse is a single entry in APICasesResponse
type APICasesResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Cases     int       `json:"CASES"`
}

// APICasesResponses is a slice of APICasesResponses structs
//easyjson:json
type APICasesResponses []*APICasesResponse

var _ apiclient.APIResponse = &APICasesResponse{}

// GetTimestamp returns the entry's timestamp
func (v APICasesResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APICasesResponse) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case apiclient.GroupByRegion:
		value = v.Region
	case apiclient.GroupByProvince:
		value = v.Province
	case apiclient.GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APICasesResponse) GetTotalValue() float64 {
	return float64(v.Cases)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APICasesResponse) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APICasesResponse) GetAttributeValues() (values []float64) {
	return []float64{float64(v.Cases)}
}

// GetCases retrieves all recorded COVID-19 cases
func (client *Client) GetCases(ctx context.Context) (results []apiclient.APIResponse, err error) {
	var body []byte
	body, err = client.call(ctx, "cases")
	if err != nil {
		return
	}

	var response APICasesResponses
	err = easyjson.Unmarshal(body, &response)
	if err != nil {
		return
	}

	return copyMaybeSort(response), nil
}
