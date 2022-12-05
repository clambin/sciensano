package simplejsonserver

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/clambin/simplejson/v5"
)

type Handler struct {
	Fetch      func(context.Context, sciensano.SummaryColumn) (*tabulator.Tabulator, error)
	Mode       sciensano.SummaryColumn
	Accumulate bool
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (h Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: h.query}
}

func (h Handler) query(ctx context.Context, req simplejson.QueryRequest) (simplejson.Response, error) {
	records, err := h.Fetch(ctx, h.Mode)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	records = records.Copy()
	if h.Accumulate {
		records.Accumulate()
	}
	records.Filter(req.Args.Range.From, req.Args.Range.To)
	return createTableResponse(records), nil
}

type Handler2 struct {
	Fetch      func(context.Context, sciensano.SummaryColumn, sciensano.DoseType) (*tabulator.Tabulator, error)
	Mode       sciensano.SummaryColumn
	DoseType   sciensano.DoseType
	Accumulate bool
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (h Handler2) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: h.query}
}

func (h Handler2) query(ctx context.Context, req simplejson.QueryRequest) (simplejson.Response, error) {
	records, err := h.Fetch(ctx, h.Mode, h.DoseType)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	records = records.Copy()
	if h.Accumulate {
		records.Accumulate()
	}
	records.Filter(req.Args.Range.From, req.Args.Range.To)
	return createTableResponse(records), nil
}

func createTableResponse(t *tabulator.Tabulator) simplejson.Response {
	columnNames := t.GetColumns()
	columns := make([]simplejson.Column, 1+len(columnNames))
	columns[0] = simplejson.Column{Text: "time", Data: simplejson.TimeColumn(t.GetTimestamps())}
	for index, column := range t.GetColumns() {
		values, _ := t.GetValues(column)
		columns[index+1] = simplejson.Column{Text: column, Data: simplejson.NumberColumn(values)}
	}

	return simplejson.TableResponse{Columns: columns}
}
