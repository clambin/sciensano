package server

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/cache/sciensano"
)

var _ grafanaJSONServer.Handler = Handler{}

type Handler struct {
	Fetch func(context.Context, sciensano.SummaryColumn) (*tabulator.Tabulator, error)
	grafanaJSONServer.Metric
	Accumulate bool
}

func (h Handler) Query(ctx context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
	mode, err := getSummaryMode(target, request)
	if err != nil {
		return nil, err
	}
	records, err := h.Fetch(ctx, mode)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	records = records.Copy()
	if h.Accumulate {
		records.Accumulate()
	}
	records.Filter(request.Range.From, request.Range.To)
	return createTableResponse(records), nil
}

var _ grafanaJSONServer.Handler = Handler2{}

type Handler2 struct {
	Fetch func(context.Context, sciensano.SummaryColumn, sciensano.DoseType) (*tabulator.Tabulator, error)
	grafanaJSONServer.Metric
	DoseType   sciensano.DoseType
	Accumulate bool
}

func (h Handler2) Query(ctx context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
	mode, err := getSummaryMode(target, request)
	if err != nil {
		return nil, err
	}
	records, err := h.Fetch(ctx, mode, h.DoseType)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
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
