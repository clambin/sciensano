package sciensano

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/mailru/easyjson"
	"io"
	"time"
)

// APITestResultsResponse is the response of the Sciensano test results API
//easyjson:json
type APITestResultsResponse []*APITestResultsResponseEntry

// APITestResultsResponseEntry is a single entry in APITestResultsResponse
//easyjon:json
type APITestResultsResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

var _ measurement.Measurement = &APITestResultsResponseEntry{}

// GetTimestamp returns the entry's timestamp
func (v APITestResultsResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APITestResultsResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case measurement.GroupByRegion:
		value = v.Region
	case measurement.GroupByProvince:
		value = v.Province
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APITestResultsResponseEntry) GetTotalValue() float64 {
	return float64(v.Total)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APITestResultsResponseEntry) GetAttributeNames() []string {
	return []string{"total", "positive"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APITestResultsResponseEntry) GetAttributeValues() (values []float64) {
	return []float64{float64(v.Total), float64(v.Positive)}
}

// GetTestResults retrieves all COVID-19 test results.
func (client *Client) GetTestResults(ctx context.Context) (results []measurement.Measurement, err error) {
	var body io.ReadCloser
	body, err = client.call(ctx, "tests")

	if err == nil {
		var cvt APITestResultsResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			results = make([]measurement.Measurement, 0, len(cvt))
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}
	return
}
