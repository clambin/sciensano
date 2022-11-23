package simplejsonserver

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/clambin/sciensano/pkg/tabulator"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

type Handler struct {
	Fetch      func(sciensano.SummaryColumn) (*tabulator.Tabulator, error)
	Mode       sciensano.SummaryColumn
	Accumulate bool
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (h Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: h.query}
}

func (h Handler) query(_ context.Context, req query.Request) (query.Response, error) {
	records, err := h.Fetch(h.Mode)
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
	Fetch      func(sciensano.SummaryColumn, sciensano.DoseType) (*tabulator.Tabulator, error)
	Mode       sciensano.SummaryColumn
	DoseType   sciensano.DoseType
	Accumulate bool
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (h Handler2) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: h.query}
}

func (h Handler2) query(_ context.Context, req query.Request) (query.Response, error) {
	records, err := h.Fetch(h.Mode, h.DoseType)
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

func createTableResponse(t *tabulator.Tabulator) query.Response {
	columns := []query.Column{
		{
			Text: "time",
			Data: query.TimeColumn(t.GetTimestamps()),
		},
	}
	for _, column := range t.GetColumns() {
		values, _ := t.GetValues(column)
		columns = append(columns, query.Column{Text: column, Data: query.NumberColumn(values)})
	}

	return query.TableResponse{Columns: columns}
}
