package vaccines

import (
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// APIBatchResponse represents one batch of vaccines
type APIBatchResponse struct {
	Date         Timestamp `json:"date"`
	Manufacturer string    `json:"manufacturer"`
	Amount       int       `json:"amount"`
}

var _ apiclient.APIResponse = &APIBatchResponse{}

// GetTimestamp returns the batch's timestamp
func (b APIBatchResponse) GetTimestamp() time.Time {
	return b.Date.Time
}

// GetGroupFieldValue returns the value of a groupable field.  Not used for APIBatchResponse.
func (b APIBatchResponse) GetGroupFieldValue(groupField int) (value string) {
	if groupField == apiclient.GroupByManufacturer {
		value = b.Manufacturer
	}
	return
}

// GetTotalValue returns the entry's total number of vaccines
func (b APIBatchResponse) GetTotalValue() float64 {
	return float64(b.Amount)
}

// GetAttributeNames returns the names of the types of vaccinations
func (b APIBatchResponse) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (b APIBatchResponse) GetAttributeValues() (values []float64) {
	return []float64{float64(b.Amount)}
}

// Timestamp representation for APIBatchResponse. Needed to unmarshal the date as received from the API
type Timestamp struct {
	time.Time
}

// UnmarshalJSON unmarshals the Timestamp in a APIBatchResponse
func (date *Timestamp) UnmarshalJSON(b []byte) (err error) {
	var timestamp time.Time
	if timestamp, err = time.Parse(`"2006-01-02"`, string(b)); err == nil {
		date.Time = timestamp
	}
	return
}

// MarshalJSON marshals a Timestamp in a APIBatchResponse
func (date Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(date.Time.Format(`"2006-01-02"`)), nil
}
