package sciensano

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/mailru/easyjson"
	"io"
	"time"
)

// APIMortalityResponse is the responder of the Reporter cases API
//easyjson:json
type APIMortalityResponse []*APIMortalityResponseEntry

// APIMortalityResponseEntry is a single entry in APIMortalityResponse
//easyjson:json
type APIMortalityResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Deaths    int       `json:"DEATHS"`
}

var _ measurement.Measurement = &APIMortalityResponseEntry{}

// GetTimestamp returns the entry's timestamp
func (v APIMortalityResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APIMortalityResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case measurement.GroupByRegion:
		value = v.Region
	case measurement.GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIMortalityResponseEntry) GetTotalValue() float64 {
	return float64(v.Deaths)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIMortalityResponseEntry) GetAttributeNames() []string {
	return []string{"total"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIMortalityResponseEntry) GetAttributeValues() (values []float64) {
	return []float64{float64(v.Deaths)}
}

// GetMortality retrieves all recorded COVID-19 mortality figures
func (client *Client) GetMortality(ctx context.Context) (results []measurement.Measurement, err error) {
	var body io.ReadCloser
	body, err = client.call(ctx, "mortality")

	if err == nil {
		var cvt APIMortalityResponse
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
