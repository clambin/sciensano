package apiclient

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"io"
)

// APIHospitalisationsResponse is the response of the Sciensano cases API
//easyjson:json
type APIHospitalisationsResponse []APIHospitalisationsResponseEntry

// APIHospitalisationsResponseEntry is a single entry in APIHospitalisationsResponse
//easyjson:json
type APIHospitalisationsResponseEntry struct {
	TimeStamp   TimeStamp `json:"DATE"`
	Province    string    `json:"PROVINCE"`
	Region      string    `json:"REGION"`
	TotalIn     int       `json:"TOTAL_IN"`
	TotalInICU  int       `json:"TOTAL_IN_ICU"`
	TotalInResp int       `json:"TOTAL_IN_RESP"`
	TotalInECMO int       `json:"TOTAL_IN_ECMO"`
}

// GetHospitalisations retrieves all recorded COVID-19 cases
func (client *Client) GetHospitalisations(ctx context.Context) (results APIHospitalisationsResponse, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("hospitalisations"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_HOSP.json"); err == nil {
		err = easyjson.UnmarshalFromReader(body, &results)
		_ = body.Close()
	}

	timer.ObserveDuration()
	metricRequestsTotal.WithLabelValues("hospitalisations").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("hospitalisations").Add(1.0)
	}
	return
}
