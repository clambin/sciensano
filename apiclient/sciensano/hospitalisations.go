package sciensano

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/mailru/easyjson"
	"io"
	"time"
)

// APIHospitalisationsResponse is the responder of the Reporter hospitalisations API
//easyjson:json
type APIHospitalisationsResponse []*APIHospitalisationsResponseEntry

// APIHospitalisationsResponseEntry is a single entry in APIHospitalisationsResponse
//easyjson:json
type APIHospitalisationsResponseEntry struct {
	TimeStamp   TimeStamp `json:"DATE"`
	Province    string    `json:"PROVINCE"`
	Region      string    `json:"REGION"`
	TotalIn     int       `json:"TOTAL_IN"`
	TotalInICU  int       `json:"TOTAL_IN_ICU"`
	TotalInResp int       `json:"TOTAL_IN_RESP"`
	TotalInECMO int       `json:"TOTAL_IN_ECMO"`
}

// GetTimestamp returns the entry's timestamp
func (v APIHospitalisationsResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v APIHospitalisationsResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case measurement.GroupByRegion:
		value = v.Region
	case measurement.GroupByProvince:
		value = v.Province
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIHospitalisationsResponseEntry) GetTotalValue() float64 {
	return float64(v.TotalIn)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIHospitalisationsResponseEntry) GetAttributeNames() []string {
	return []string{"in", "inICU", "inResp", "inECMO"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIHospitalisationsResponseEntry) GetAttributeValues() (values []float64) {
	return []float64{float64(v.TotalIn), float64(v.TotalInICU), float64(v.TotalInResp), float64(v.TotalInECMO)}
}

var _ measurement.Measurement = &APIHospitalisationsResponseEntry{}

// GetHospitalisations retrieves all recorded COVID-19 cases
func (client *Client) GetHospitalisations(ctx context.Context) (results []measurement.Measurement, err error) {
	var body io.ReadCloser
	body, err = client.call(ctx, "hospitalisations")

	if err == nil {
		var cvt APIHospitalisationsResponse
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
