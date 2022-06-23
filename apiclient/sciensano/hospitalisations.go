package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// APIHospitalisationsResponse is a single entry in APIHospitalisationsResponse
type APIHospitalisationsResponse struct {
	TimeStamp   TimeStamp `json:"DATE"`
	Province    string    `json:"PROVINCE"`
	Region      string    `json:"REGION"`
	TotalIn     int       `json:"TOTAL_IN"`
	TotalInICU  int       `json:"TOTAL_IN_ICU"`
	TotalInResp int       `json:"TOTAL_IN_RESP"`
	TotalInECMO int       `json:"TOTAL_IN_ECMO"`
}

var _ apiclient.APIResponse = &APIHospitalisationsResponse{}

// GetTimestamp returns the entry's timestamp
func (v APIHospitalisationsResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APIHospitalisationsResponse) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case apiclient.GroupByRegion:
		value = v.Region
	case apiclient.GroupByProvince:
		value = v.Province
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIHospitalisationsResponse) GetTotalValue() float64 {
	return float64(v.TotalIn)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIHospitalisationsResponse) GetAttributeNames() []string {
	return []string{"in", "inICU", "inResp", "inECMO"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIHospitalisationsResponse) GetAttributeValues() (values []float64) {
	return []float64{float64(v.TotalIn), float64(v.TotalInICU), float64(v.TotalInResp), float64(v.TotalInECMO)}
}
