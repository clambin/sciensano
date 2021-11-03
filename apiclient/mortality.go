package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"io"
)

// APIMortalityResponse is the response of the Sciensano cases API
//easyjson:json
type APIMortalityResponse []APIMortalityResponseEntry

// APIMortalityResponseEntry is a single entry in APIMortalityResponse
//easyjson:json
type APIMortalityResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Deaths    int       `json:"DEATHS"`
}

// GetMortality retrieves all recorded COVID-19 mortality figures
func (client *Client) GetMortality(ctx context.Context) (results APIMortalityResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("mortality"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_MORT.json"); err == nil {
		err = easyjson.UnmarshalFromReader(body, &results)
	}
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("mortality").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("mortality").Add(1.0)
	}
	return
}
