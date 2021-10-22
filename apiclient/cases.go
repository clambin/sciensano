package apiclient

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

type APICasesResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Cases     int       `json:"CASES"`
}

func (client *Client) GetCases(ctx context.Context) (results []*APICasesResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("cases"))
	err = client.call(ctx, "COVID19BE_CASES_AGESEX.json", &results)
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("cases").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("cases").Add(1.0)
	}
	return
}
