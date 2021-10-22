package apiclient

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

type APIVaccinationsResponse struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Gender    string    `json:"SEX"`
	Dose      string    `json:"DOSE"`
	Count     int       `json:"Count"`
}

// GetVaccinations retrieves all COVID-19 vaccinations.
func (client *Client) GetVaccinations(ctx context.Context) (results []*APIVaccinationsResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("vaccinations"))
	err = client.call(ctx, "COVID19BE_VACC.json", &results)
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("vaccinations").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("vaccinations").Add(1.0)
	}
	return
}
