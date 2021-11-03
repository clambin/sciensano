package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"io"
)

// APIVaccinationsResponse is the response of the Sciensano vaccinations API
//easyjson:json
type APIVaccinationsResponse []APIVaccinationsResponseEntry

// APIVaccinationsResponseEntry is a single entry in APIVaccinationResponse
//easyjson:json
type APIVaccinationsResponseEntry struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Gender    string    `json:"SEX"`
	Dose      string    `json:"DOSE"`
	Count     int       `json:"COUNT"`
}

// GetVaccinations retrieves all COVID-19 vaccinations.
func (client *Client) GetVaccinations(ctx context.Context) (results APIVaccinationsResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("vaccinations"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_VACC.json"); err == nil {
		err = easyjson.UnmarshalFromReader(body, &results)
		_ = body.Close()
	}
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("vaccinations").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("vaccinations").Add(1.0)
	}
	return
}
