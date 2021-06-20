package apihandler

import (
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"time"
)

func (handler *Handler) buildTestTableResponse(_, endTime time.Time, _ string) (response *grafana_json.TableQueryResponse) {
	if tests, err := handler.Sciensano.GetTests(endTime); err == nil {

		rows := len(tests)
		timestamps := make(grafana_json.TableQueryResponseTimeColumn, rows)
		allTests := make(grafana_json.TableQueryResponseNumberColumn, rows)
		positiveTests := make(grafana_json.TableQueryResponseNumberColumn, rows)
		positiveRate := make(grafana_json.TableQueryResponseNumberColumn, rows)

		for index, test := range tests {
			timestamps[index] = test.Timestamp
			allTests[index] = float64(test.Total)
			positiveTests[index] = float64(test.Positive)
			if test.Total != 0 {
				positiveRate[index] = float64(test.Positive) / float64(test.Total)
			} else {
				positiveRate[index] = 0
			}
		}

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "total", Data: allTests},
			{Text: "positive", Data: positiveTests},
			{Text: "rate", Data: positiveRate},
		}
	}
	return
}

func (handler *Handler) buildTestForecastTableResponse(_, endTime time.Time, _ string) (response *grafana_json.TableQueryResponse) {
	if tests, err := handler.Sciensano.GetTests(endTime); err == nil {

		var forecast []sciensano.Test
		forecast, err = predictor.ForecastTests(tests)

		if err != nil {
			log.WithError(err).Warning("unable to forecast tests")
			return nil
		}

		rows := len(forecast)
		timestamps := make(grafana_json.TableQueryResponseTimeColumn, rows)
		totalTests := make(grafana_json.TableQueryResponseNumberColumn, rows)
		positiveTests := make(grafana_json.TableQueryResponseNumberColumn, rows)
		positiveRate := make(grafana_json.TableQueryResponseNumberColumn, rows)

		for index, test := range forecast {
			timestamps[index] = test.Timestamp
			totalTests[index] = float64(test.Total)
			positiveTests[index] = float64(test.Positive)
			if test.Total != 0 {
				positiveRate[index] = float64(test.Positive) / float64(test.Total)
			} else {
				positiveRate[index] = 0
			}
		}

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "total", Data: totalTests},
			{Text: "positive", Data: positiveTests},
			{Text: "rate", Data: positiveRate},
		}
	}
	return
}
