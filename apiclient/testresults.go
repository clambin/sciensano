package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"io"
)

// APITestResultsResponse is the response of the Sciensano test results API
//easyjson:json
type APITestResultsResponse []APITestResultsResponseEntry

// APITestResultsResponseEntry is a single entry in APITestResultsResponse
//easyjon:json
type APITestResultsResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

// GetTestResults retrieves all COVID-19 test results.
func (client *Client) GetTestResults(ctx context.Context) (results APITestResultsResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("tests"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_tests.json"); err == nil {
		err = easyjson.UnmarshalFromReader(body, &results)
		_ = body.Close()
	}
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("tests").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("tests").Add(1.0)
	}
	return
}
