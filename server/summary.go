package server

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
)

var _ grafanaJSONServer.Handler = SummaryHandler{}

//go:generate mockery --name ReportsStore --with-expecter=true
type ReportsStore interface {
	Get(string) (*tabulator.Tabulator, error)
}

type SummaryHandler struct {
	ReportsStore
	grafanaJSONServer.Metric
	Accumulate bool // TODO: here or in reporter?
	metricProcessor
}

func (h SummaryHandler) Query(_ context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
	key, err := h.metricProcessor.getKey(target, request)
	if err != nil {
		return nil, fmt.Errorf("unable to get store key: %w", err)
	}

	records, err := h.ReportsStore.Get(key)
	if err != nil {
		return nil, fmt.Errorf("fetch %s failed: %w", key, err)
	}
	records = records.Copy()
	if h.Accumulate {
		records.Accumulate()
	}
	records.Filter(request.Range.From, request.Range.To)
	return createTableResponse(records), nil
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
