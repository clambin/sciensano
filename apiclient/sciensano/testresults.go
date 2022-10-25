package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// APITestResultsResponse is a single entry in APITestResultsResponse
type APITestResultsResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

var _ apiclient.APIResponse = &APITestResultsResponse{}

// GetTimestamp returns the entry's timestamp
func (v APITestResultsResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APITestResultsResponse) GetGroupFieldValue(groupField apiclient.GroupField) (value string) {
	switch groupField {
	case apiclient.GroupByRegion:
		value = v.Region
	case apiclient.GroupByProvince:
		value = v.Province
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APITestResultsResponse) GetTotalValue() float64 {
	return float64(v.Total)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APITestResultsResponse) GetAttributeNames() []string {
	return []string{"total", "positive"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APITestResultsResponse) GetAttributeValues() (values []float64) {
	return []float64{float64(v.Total), float64(v.Positive)}
}
