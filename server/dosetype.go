package server

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/set"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
)

type SummaryByDoseTypeHandler struct {
	ReportsStore
	grafanaJSONServer.Metric
}

func newSummaryByDoseTypeHandler(name string, summaryColumns set.Set[sciensano.SummaryColumn], doseTypes set.Set[sciensano.DoseType], s ReportsStore) SummaryByDoseTypeHandler {
	var summaryPayloadOptions []grafanaJSONServer.MetricPayloadOption
	var doseTypePayloadOptions []grafanaJSONServer.MetricPayloadOption

	for _, summaryColumn := range summaryColumns.List() {
		summaryPayloadOptions = append(summaryPayloadOptions, grafanaJSONServer.MetricPayloadOption{
			Label: summaryColumn.String(),
			Value: summaryColumn.String(),
		})
	}
	for _, doseType := range doseTypes.List() {
		doseTypePayloadOptions = append(summaryPayloadOptions, grafanaJSONServer.MetricPayloadOption{
			Label: doseType.String(),
			Value: doseType.String(),
		})
	}
	return SummaryByDoseTypeHandler{
		ReportsStore: s,
		Metric: grafanaJSONServer.Metric{Value: name, Payloads: []grafanaJSONServer.MetricPayload{
			{
				Label:   "Summary",
				Name:    "summary",
				Type:    "select",
				Width:   40,
				Options: summaryPayloadOptions,
			},
			{
				Label:   "DoseType",
				Name:    "doseType",
				Type:    "select",
				Width:   40,
				Options: doseTypePayloadOptions,
			},
			{
				Label: "Accumulate",
				Name:  "accumulate",
				Type:  "select",
				Width: 40,
				Options: []grafanaJSONServer.MetricPayloadOption{
					{Label: "Yes", Value: "yes"},
					{Label: "No", Value: "no"},
				},
			},
		}},
	}
}

func (h SummaryByDoseTypeHandler) Query(_ context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
	key, accumulate, err := h.getRequestOptions(target, request)
	if err != nil {
		return nil, fmt.Errorf("unable to get store key: %w", err)
	}

	records, err := h.ReportsStore.Get(key)
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

func (h SummaryByDoseTypeHandler) getRequestOptions(target string, req grafanaJSONServer.QueryRequest) (string, bool, error) {
	var summaryOption struct {
		Summary    string
		DoseType   string
		Accumulate string
	}
	var accumulate bool
	if err := req.GetPayload(target, &summaryOption); err != nil {
		return "", accumulate, fmt.Errorf("invalid payload: %w", err)
	}

	//slog.Debug("getting request options", "row", string(req.Targets[0].Payload), "options", summaryOption)

	//mode, ok := sciensano.SummaryColumnNames[summaryOption.Summary]
	//if !ok {
	//	return "", accumulate, fmt.Errorf("invalid summary option: %s", summaryOption.Summary)
	//}
	switch summaryOption.Accumulate {
	case "yes":
		accumulate = true
	case "no":
		accumulate = false
	default:
		return "", accumulate, fmt.Errorf("invalid accumulate value: %s", summaryOption.Accumulate)
	}

	return target + "-" + summaryOption.Summary + summaryOption.DoseType, accumulate, nil
}
