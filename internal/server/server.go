package server

import (
	"context"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
	gjson "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/v2/internal/sciensano"
	"log/slog"
)

// Server groups Grafana JSON API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	JSONServer *gjson.Server
	Handlers   map[string]gjson.Handler
	reports    ReportsStore
}

type ReportsStore interface {
	Get(string) (*tabulator.Tabulator, error)
	Keys() []string
}

func New(reportsStore ReportsStore, metrics gjson.PrometheusQueryMetrics, logger *slog.Logger) *Server {
	s := &Server{
		Handlers: make(map[string]gjson.Handler),
		reports:  reportsStore,
	}

	options := []gjson.Option{
		gjson.WithLogger(logger.With("component", "grafana-json-server")),
		gjson.WithPrometheusQueryMetrics(metrics),
	}

	summaryHandlers := []struct {
		name           string
		summaryColumns set.Set[sciensano.SummaryColumn]
		accumulate     bool
	}{
		{name: "cases", summaryColumns: sciensano.CasesValidSummaryModes()},
		{name: "hospitalisations", summaryColumns: sciensano.HospitalisationsValidSummaryModes()},
		{name: "mortalities", summaryColumns: sciensano.MortalitiesValidSummaryModes()},
		{name: "tests", summaryColumns: sciensano.TestResultsValidSummaryModes()},
		{name: "vaccinations", summaryColumns: sciensano.VaccinationsValidSummaryModes(), accumulate: true},
	}

	for _, summaryHandler := range summaryHandlers {
		metric, h := newSummaryMetric(reportsStore, summaryHandler.name, summaryHandler.summaryColumns.List())

		s.Handlers[summaryHandler.name] = h
		options = append(options, gjson.WithMetric(metric, h, nil))
	}

	metric, h := newVaccinationDoseTypeMetric(reportsStore, "vaccination-rate", []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup}, []sciensano.DoseType{sciensano.Partial, sciensano.Full})
	s.Handlers[metric.Value] = h
	options = append(options, gjson.WithMetric(metric, h, nil))

	s.JSONServer = gjson.NewServer(options...)
	s.JSONServer.HandleFunc("/health", s.Health)
	return s
}

// Run starts the supporting components
func (s *Server) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func createTableResponse(t *tabulator.Tabulator) gjson.QueryResponse {
	columnNames := t.GetColumns()
	columns := make([]gjson.Column, 1+len(columnNames))
	columns[0] = gjson.Column{Text: "time", Data: gjson.TimeColumn(t.GetTimestamps())}
	for index, column := range t.GetColumns() {
		values, _ := t.GetValues(column)
		columns[index+1] = gjson.Column{Text: column, Data: gjson.NumberColumn(values)}
	}

	return gjson.TableResponse{Columns: columns}
}
