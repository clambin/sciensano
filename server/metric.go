package server

import (
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"strconv"
)

func newMetric(name string, summaryColumns ...sciensano.SummaryColumn) grafanaJSONServer.Metric {
	var options []grafanaJSONServer.MetricPayloadOption

	for _, summaryColumn := range summaryColumns {
		options = append(options, grafanaJSONServer.MetricPayloadOption{
			Label: summaryColumn.String(),
			Value: strconv.Itoa(int(summaryColumn)),
		})
	}
	return grafanaJSONServer.Metric{Value: name, Payloads: []grafanaJSONServer.MetricPayload{{
		Label:   "Summary",
		Name:    "summary",
		Type:    "select",
		Width:   40,
		Options: options,
	}}}
}

func getSummaryMode(target string, req grafanaJSONServer.QueryRequest) (sciensano.SummaryColumn, error) {
	panic("not implemented")
}
