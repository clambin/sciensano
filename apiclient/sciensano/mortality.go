package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// APIMortalityResponse is a single entry in APIMortalityResponse
type APIMortalityResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Deaths    int       `json:"DEATHS"`
}

var _ apiclient.APIResponse = &APIMortalityResponse{}

// GetTimestamp returns the entry's timestamp
func (v APIMortalityResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APIMortalityResponse) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case apiclient.GroupByRegion:
		value = v.Region
	case apiclient.GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIMortalityResponse) GetTotalValue() float64 {
	return float64(v.Deaths)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIMortalityResponse) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIMortalityResponse) GetAttributeValues() (values []float64) {
	return []float64{float64(v.Deaths)}
}
