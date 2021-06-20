package apihandler

import (
	grafana_json "github.com/clambin/grafana-json"
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
			positiveRate[index] = float64(test.Positive) / float64(test.Total)
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
