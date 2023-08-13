package server

import (
	"context"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/go-common/set"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
)

// Server groups Grafana JSON API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	JSONServer *grafanaJSONServer.Server
	Handlers   map[string]grafanaJSONServer.Handler
}

var _ prometheus.Collector = &Server{}

func New(reportsStore ReportsStore, logger *slog.Logger) *Server {
	s := &Server{
		Handlers: make(map[string]grafanaJSONServer.Handler),
	}

	options := []grafanaJSONServer.Option{
		grafanaJSONServer.WithLogger(logger.With("component", "grafana-json-server")),
		grafanaJSONServer.WithRequestLogger(slog.LevelDebug, middleware.DefaultRequestLogFormatter),
		grafanaJSONServer.WithHTTPHandler(http.MethodGet, "/health", http.HandlerFunc(s.Health)),
		grafanaJSONServer.WithPrometheusQueryMetrics("sciensano", "", "sciensano"),
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
		options = append(options, grafanaJSONServer.WithMetric(metric, h, nil))
	}

	metric, h := newVaccinationDoseTypeMetric(reportsStore, "vaccination-rate", []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup}, []sciensano.DoseType{sciensano.Partial, sciensano.Full})
	s.Handlers[metric.Value] = h
	options = append(options, grafanaJSONServer.WithMetric(metric, h, nil))

	s.JSONServer = grafanaJSONServer.NewServer(options...)
	return s
}

// Run starts the supporting components
func (s *Server) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *Server) Describe(descs chan<- *prometheus.Desc) {
	// s.apiCache.Describe(descs)
	s.JSONServer.Describe(descs)
}

func (s *Server) Collect(metrics chan<- prometheus.Metric) {
	//s.apiCache.Collect(metrics)
	s.JSONServer.Collect(metrics)
}
