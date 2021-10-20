package apiclient

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
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
	results, err = client.getVaccinations(ctx)
	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("vaccinations").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("vaccinations").Add(1.0)
	}
	return
}

func (client *Client) getVaccinations(ctx context.Context) (results []*APIVaccinationsResponse, err error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.getURL()+"/Data/COVID19BE_VACC.json", nil)

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
