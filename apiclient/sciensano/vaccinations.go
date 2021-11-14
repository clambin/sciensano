package sciensano

import (
	"context"
	"github.com/clambin/sciensano/measurement"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

// APIVaccinationsResponse is the response of the Sciensano vaccinations API
//easyjson:json
type APIVaccinationsResponse []*APIVaccinationsResponseEntry

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

// GetTimestamp returns the entry's timestamp
func (v *APIVaccinationsResponseEntry) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APIVaccinationsResponseEntry) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case measurement.GroupByRegion:
		value = v.Region
	case measurement.GroupByAgeGroup:
		value = v.AgeGroup
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIVaccinationsResponseEntry) GetTotalValue() float64 {
	return float64(v.Count)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIVaccinationsResponseEntry) GetAttributeNames() []string {
	return []string{"partial", "full", "singledose", "booster"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIVaccinationsResponseEntry) GetAttributeValues() (values []float64) {
	values = []float64{0, 0, 0, 0}
	switch v.Dose {
	case "A":
		values[0] = float64(v.Count)
	case "B":
		values[1] = float64(v.Count)
	case "C":
		values[2] = float64(v.Count)
	case "E":
		values[3] = float64(v.Count)
	}
	return
}

var _ measurement.Measurement = &APIVaccinationsResponseEntry{}

// GetVaccinations retrieves all COVID-19 vaccinations.
func (client *Client) GetVaccinations(ctx context.Context) (results []measurement.Measurement, err error) {
	return client.Call(ctx, "vaccinations", client.getVaccinations)
}

func (client *Client) getVaccinations(ctx context.Context) (results []measurement.Measurement, err error) {
	timer := prometheus.NewTimer(metricRequestLatency.WithLabelValues("vaccinations"))
	var body io.ReadCloser
	if body, err = client.call(ctx, "COVID19BE_VACC.json"); err == nil {
		var cvt APIVaccinationsResponse
		if err = easyjson.UnmarshalFromReader(body, &cvt); err == nil {
			results = make([]measurement.Measurement, 0, len(cvt))
			for _, entry := range cvt {
				results = append(results, entry)
			}
		}
		_ = body.Close()
	}
	duration := timer.ObserveDuration()
	log.WithField("duration", duration).Debug("called GetVaccinations API")
	metricRequestsTotal.WithLabelValues("vaccinations").Add(1.0)
	if err != nil {
		metricRequestErrorsTotal.WithLabelValues("vaccinations").Add(1.0)
	}
	return
}
