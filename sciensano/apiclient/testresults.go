package apiclient

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
)

type APITestResultsResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

func (client *Client) GetTestResults(ctx context.Context) (results []*APITestResultsResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("tests"))
	results, err = client.getTestResults(ctx)
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("tests").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("tests").Add(1.0)
	}
	return
}

func (client *Client) getTestResults(ctx context.Context) (results []*APITestResultsResponse, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.getURL()+"/Data/COVID19BE_tests.json", nil)

	var resp *http.Response
	resp, err = client.HTTPClient.Do(req)

	if err != nil {
		return
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	var body []byte
	body, _ = io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &results)

	return
}
