package sciensano

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/mailru/easyjson"
	"io"
	"time"
)

// APICasesResponse is the response of the Sciensano cases API
//easyjson:json
type APICasesResponse []*APICasesResponseEntry

// APICasesResponseEntry is a single entry in APICasesResponse
//easyjson:json
type APICasesResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Cases     int       `json:"CASES"`
}

var _ measurement.Measurement = &APICasesResponseEntry{}

// GetTimestamp returns the entry's timestamp
func (v APICasesResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APICasesResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case measurement.GroupByRegion:
		value = v.Region
	case measurement.GroupByProvince:
		value = v.Province
	case measurement.GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APICasesResponseEntry) GetTotalValue() float64 {
	return float64(v.Cases)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APICasesResponseEntry) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APICasesResponseEntry) GetAttributeValues() (values []float64) {
	return []float64{float64(v.Cases)}
}

// GetCases retrieves all recorded COVID-19 cases
func (client *Client) GetCases(ctx context.Context) (results []measurement.Measurement, err error) {
	var body io.ReadCloser
	body, err = client.call(ctx, "cases")

	if err == nil {
		var cvt APICasesResponse
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
