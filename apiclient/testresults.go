package apiclient

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

// APITestResultsResponse is the response of the Sciensano test results API
type APITestResultsResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

// GetTestResults retrieves all COVID-19 test results.
func (client *Client) GetTestResults(ctx context.Context) (results []*APITestResultsResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("tests"))
	err = client.call(ctx, "COVID19BE_tests.json", &results)
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("tests").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("tests").Add(1.0)
	}
	return
}
