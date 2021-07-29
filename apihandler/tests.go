package apihandler

import (
	grafanaJson "github.com/clambin/grafana-json"
	"time"
)

func (handler *Handler) buildTestTableResponse(_, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse) {
	if tests, err := handler.Sciensano.GetTests(endTime); err == nil {

		rows := len(tests)
		timestamps := make(grafanaJson.TableQueryResponseTimeColumn, rows)
		allTests := make(grafanaJson.TableQueryResponseNumberColumn, rows)
		positiveTests := make(grafanaJson.TableQueryResponseNumberColumn, rows)
		positiveRate := make(grafanaJson.TableQueryResponseNumberColumn, rows)

		for index, test := range tests {
			timestamps[index] = test.Timestamp
			allTests[index] = float64(test.Total)
			positiveTests[index] = float64(test.Positive)
			positiveRate[index] = float64(test.Positive) / float64(test.Total)
		}

		response = new(grafanaJson.TableQueryResponse)
		response.Columns = []grafanaJson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "total", Data: allTests},
			{Text: "positive", Data: positiveTests},
			{Text: "rate", Data: positiveRate},
		}
	}
	return
}
