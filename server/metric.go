package server

import (
	"context"
	"fmt"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
)

func newSummaryMetric(s ReportsStore, name string, values []sciensano.SummaryColumn) (grafanaJSONServer.Metric, grafanaJSONServer.Handler) {
	var v []fmt.Stringer
	for _, value := range values {
		v = append(v, value)
	}
	metric := makeMetric(name, []metricOption{{name: "Summary", values: v}}...)
	return metric, handler{s: s, parseRequest: parseSummaryRequest}
}

type metricOption struct {
	name   string
	values []fmt.Stringer
}

func makeMetric(name string, options ...metricOption) grafanaJSONServer.Metric {
	var payloads []grafanaJSONServer.MetricPayload
	for _, option := range options {
		var payloadOptions []grafanaJSONServer.MetricPayloadOption
		for _, value := range option.values {
			payloadOptions = append(payloadOptions, grafanaJSONServer.MetricPayloadOption{
				Label: value.String(),
				Value: value.String(),
			})
		}
		payloads = append(payloads, grafanaJSONServer.MetricPayload{
			Label: option.name,
			Name:  option.name,
			Type:  "select",
			// ReloadMetric: false,
			Width:   40,
			Options: payloadOptions,
		})
	}
	payloads = append(payloads, grafanaJSONServer.MetricPayload{
		Label: "Accumulate",
		Name:  "Accumulate",
		Type:  "select",
		Width: 40,
		Options: []grafanaJSONServer.MetricPayloadOption{
			{Label: "Yes", Value: "yes"},
			{Label: "No", Value: "no"},
		},
	})

	return grafanaJSONServer.Metric{Value: name, Label: name, Payloads: payloads}
}

type handler struct {
	s            ReportsStore
	parseRequest func(string, grafanaJSONServer.QueryRequest) (string, bool, error)
}

func (h handler) Query(_ context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
	key, accumulate, err := h.parseRequest(target, request)
	if err != nil {
		return nil, fmt.Errorf("unable to get store key: %w", err)
	}

	records, err := h.s.Get(key)
	if err != nil {
		return nil, fmt.Errorf("fetch %s failed: %w", key, err)
	}
	records = records.Copy()
	if accumulate {
		records.Accumulate()
	}
	records.Filter(request.Range.From, request.Range.To)
	return createTableResponse(records), nil
}

func parseSummaryRequest(target string, req grafanaJSONServer.QueryRequest) (string, bool, error) {
	var summaryOption struct {
		Summary    string
		Accumulate string
	}
	var accumulate bool
	if err := req.GetPayload(target, &summaryOption); err != nil {
		return "", accumulate, fmt.Errorf("invalid payload: %w", err)
	}

	//slog.Debug("getting request options", "row", string(req.Targets[0].Payload), "options", summaryOption)

	mode, ok := sciensano.SummaryColumnNames[summaryOption.Summary]
	if !ok {
		return "", accumulate, fmt.Errorf("invalid summary option: %s", summaryOption.Summary)
	}
	switch summaryOption.Accumulate {
	case "yes":
		accumulate = true
	case "no":
		accumulate = false
	default:
		return "", accumulate, fmt.Errorf("invalid accumulate value: %s", summaryOption.Accumulate)
	}

	return target + "-" + mode.String(), accumulate, nil
}
