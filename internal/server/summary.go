package server

import (
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
)

//go:generate mockery --name ReportsStore --with-expecter=true
type ReportsStore interface {
	Get(string) (*tabulator.Tabulator, error)
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
