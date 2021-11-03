package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"io"
)

// APICasesResponse is the response of the Sciensano cases API
//easyjson:json
type APICasesResponse []APICasesResponseEntry

// APICasesResponseEntry is a single entry in APICasesResponse
//easyjson:json
type APICasesResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Cases     int       `json:"CASES"`
}

// GetCases retrieves all recorded COVID-19 cases
func (client *Client) GetCases(ctx context.Context) (results APICasesResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("cases"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_CASES_AGESEX.json"); err == nil {
		err = easyjson.UnmarshalFromReader(body, &results)
		_ = body.Close()
	}

	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("cases").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("cases").Add(1.0)
	}
	return
}
