package server

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"golang.org/x/exp/slog"
)

var _ grafanaJSONServer.Handler = SummaryHandler{}

//go:generate mockery --name ReportsStore --with-expecter=true
type ReportsStore interface {
	Get(string) (*tabulator.Tabulator, error)
}

type SummaryHandler struct {
	ReportsStore
	grafanaJSONServer.Metric
}

func newSummaryHandler(name string, summaryColumns set.Set[sciensano.SummaryColumn], s ReportsStore) SummaryHandler {
	var summaryPayloadOptions []grafanaJSONServer.MetricPayloadOption

	for _, summaryColumn := range summaryColumns.List() {
		summaryPayloadOptions = append(summaryPayloadOptions, grafanaJSONServer.MetricPayloadOption{
			Label: summaryColumn.String(),
			Value: summaryColumn.String(),
		})
	}
	return SummaryHandler{
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

func (h SummaryHandler) Query(_ context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
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

func (h SummaryHandler) getRequestOptions(target string, req grafanaJSONServer.QueryRequest) (string, bool, error) {
	var summaryOption struct {
		Summary    string
		Accumulate string
	}
	var accumulate bool
	if err := req.GetPayload(target, &summaryOption); err != nil {
		return "", accumulate, fmt.Errorf("invalid payload: %w", err)
	}

	slog.Debug("getting request options", "row", string(req.Targets[0].Payload), "options", summaryOption)

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

func createTableResponse(t *tabulator.Tabulator) grafanaJSONServer.QueryResponse {
	columnNames := t.GetColumns()
	columns := make([]grafanaJSONServer.Column, 1+len(columnNames))
	columns[0] = grafanaJSONServer.Column{Text: "time", Data: grafanaJSONServer.TimeColumn(t.GetTimestamps())}
	for index, column := range t.GetColumns() {
		values, _ := t.GetValues(column)
		columns[index+1] = grafanaJSONServer.Column{Text: column, Data: grafanaJSONServer.NumberColumn(values)}
	}

	return grafanaJSONServer.TableResponse{Columns: columns}
}
