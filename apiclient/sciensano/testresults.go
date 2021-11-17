package sciensano

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/metrics"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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
	timer := prometheus.NewTimer(metrics.MetricRequestLatency.WithLabelValues("tests"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_tests.json"); err == nil {
		var cvt APITestResultsResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			results = make([]measurement.Measurement, 0, len(cvt))
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}
	duration := timer.ObserveDuration()
	log.WithField("duration", duration).Debug("called GetTestResults API")
	metrics.MetricRequestsTotal.WithLabelValues("tests").Add(1.0)
	if err != nil {
		metrics.MetricRequestErrorsTotal.WithLabelValues("tests").Add(1.0)
	}
	return
}
