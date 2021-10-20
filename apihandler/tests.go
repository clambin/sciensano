package apihandler

import (
	"context"
	"fmt"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"time"
)

func (handler *Handler) buildTestTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse, err error) {
	var tests []sciensano.TestResult
	tests, err = handler.Sciensano.GetTests(ctx, endTime)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve test results: %s", err.Error())
	}

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

	return
}
