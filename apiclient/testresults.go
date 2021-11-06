package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
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

// GetTimestamp returns the entry's timestamp
func (v *APITestResultsResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APITestResultsResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case GroupByRegion:
		value = v.Region
	case GroupByProvince:
		value = v.Province
	}
	return
}

// GetTestResults retrieves all COVID-19 test results.
func (client *Client) GetTestResults(ctx context.Context) (results []Measurement, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("tests"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_tests.json"); err == nil {
		var cvt APITestResultsResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("tests").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("tests").Add(1.0)
	}
	return
}
