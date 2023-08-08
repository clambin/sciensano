package server

import (
	"fmt"
	"github.com/clambin/go-common/set"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
)

func newSummaryMetric(name string, summaryColumns set.Set[sciensano.SummaryColumn]) grafanaJSONServer.Metric {
	f := summaryMetricProcessor{summaryColumns: summaryColumns}

	return grafanaJSONServer.Metric{Value: name, Payloads: f.makeMetricPayload()}
}

type metricProcessor interface {
	makeMetricPayload() []grafanaJSONServer.MetricPayload
	getKey(target string, req grafanaJSONServer.QueryRequest) (string, error)
}

type summaryMetricProcessor struct {
	summaryColumns set.Set[sciensano.SummaryColumn]
}

func (f summaryMetricProcessor) makeMetricPayload() []grafanaJSONServer.MetricPayload {
	var options []grafanaJSONServer.MetricPayloadOption

	for _, summaryColumn := range f.summaryColumns.List() {
		options = append(options, grafanaJSONServer.MetricPayloadOption{
			Label: summaryColumn.String(),
			Value: summaryColumn.String(),
		})
	}
	return []grafanaJSONServer.MetricPayload{
		{
			Label:   "Summary",
			Name:    "summary",
			Type:    "select",
			Width:   40,
			Options: options,
		},
	}
}

func (f summaryMetricProcessor) getKey(target string, req grafanaJSONServer.QueryRequest) (string, error) {
	var summaryOption struct {
		Summary string
	}
	if err := req.GetPayload(target, &summaryOption); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}

	mode, ok := sciensano.SummaryColumnNames[summaryOption.Summary]
	if !ok {
		return "", fmt.Errorf("invalid summary option: %s", summaryOption.Summary)
	}
	return target + "-" + mode.String(), nil
}
