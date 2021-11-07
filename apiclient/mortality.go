package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

// APIMortalityResponse is the response of the Sciensano cases API
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

// GetTimestamp returns the entry's timestamp
func (v *APIMortalityResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APIMortalityResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case GroupByRegion:
		value = v.Region
	case GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetMortality retrieves all recorded COVID-19 mortality figures
func (client *Client) GetMortality(ctx context.Context) (results []Measurement, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("mortality"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_MORT.json"); err == nil {
		var cvt APIMortalityResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			results = make([]Measurement, 0, len(cvt))
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}
	duration := timer.ObserveDuration()
	log.WithField("duration", duration).Debug("called GetMortality API")
	metricRequestsTotal.WithLabelValues("mortality").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("mortality").Add(1.0)
	}
	return
}
