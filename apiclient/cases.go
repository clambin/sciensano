package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
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

// GetTimestamp returns the entry's timestamp
func (v *APICasesResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APICasesResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case GroupByRegion:
		value = v.Region
	case GroupByProvince:
		value = v.Province
	case GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetCases retrieves all recorded COVID-19 cases
func (client *Client) GetCases(ctx context.Context) (results []Measurement, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("cases"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_CASES_AGESEX.json"); err == nil {
		var cvt APICasesResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}

	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("cases").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("cases").Add(1.0)
	}
	return
}
