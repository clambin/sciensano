package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"time"
)

// APIVaccinationsResponse is the response of the Sciensano vaccinations API
//easyjson:json
type APIVaccinationsResponse []*APIVaccinationsResponseEntry

// APIVaccinationsResponseEntry is a single entry in APIVaccinationResponse
//easyjson:json
type APIVaccinationsResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Gender    string    `json:"SEX"`
	Dose      string    `json:"DOSE"`
	Count     int       `json:"COUNT"`
}

// GetTimestamp returns the entry's timestamp
func (v *APIVaccinationsResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APIVaccinationsResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case GroupByRegion:
		value = v.Region
	case GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetVaccinations retrieves all COVID-19 vaccinations.
func (client *Client) GetVaccinations(ctx context.Context) (results []Measurement, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("vaccinations"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_VACC.json"); err == nil {
		var cvt APIVaccinationsResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			results = make([]Measurement, 0, len(cvt))
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("vaccinations").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("vaccinations").Add(1.0)
	}
	return
}
