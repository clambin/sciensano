package queryhandlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"strconv"
)

var _ grafanaJSONServer.Handler = SummaryHandler{}

//go:generate mockery --name ReportsStore --with-expecter=true
type ReportsStore interface {
	Get(string) (*tabulator.Tabulator, error)
}

type SummaryHandler struct {
	ReportsStore
	grafanaJSONServer.Metric
	Accumulate bool
}

func (h SummaryHandler) Query(_ context.Context, target string, request grafanaJSONServer.QueryRequest) (grafanaJSONServer.QueryResponse, error) {
	mode, err := getSummaryMode(target, request)
	if err != nil {
		return nil, fmt.Errorf("summary mode: %w", err)
	}
	key := target + "-" + mode.String()

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

func getSummaryMode(target string, req grafanaJSONServer.QueryRequest) (sciensano.SummaryColumn, error) {
	var summaryOption struct {
		Summary string
	}
	var summary int
	err := req.GetPayload(target, &summaryOption)
	if err == nil {
		summary, err = strconv.Atoi(summaryOption.Summary)
	}
	if summary < 0 {
		err = errors.New("invalid mode")
	}
	if err != nil {
		return -1, fmt.Errorf(`invalid summary in payload "%s": %w`, summaryOption.Summary, err)
	}
	return sciensano.SummaryColumn(summary), nil
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
