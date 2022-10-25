package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
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
//
//easyjson:json
type APICasesResponses []*APICasesResponse

var _ apiclient.APIResponse = &APICasesResponse{}

// GetTimestamp returns the entry's timestamp
func (v APICasesResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APICasesResponse) GetGroupFieldValue(groupField apiclient.GroupField) (value string) {
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
