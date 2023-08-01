package server

import (
	"context"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/queryhandlers"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
)

// Server groups Grafana JSON API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	JSONServer *grafanaJSONServer.Server
	Handlers   map[string]grafanaJSONServer.Handler
}

//go:generate mockery --name ReportsStorer --with-expecter=true
type ReportsStorer interface {
	Get(string) (*tabulator.Tabulator, error)
	Put(string, *tabulator.Tabulator)
}

var _ prometheus.Collector = &Server{}

func New(reportsStore ReportsStorer, logger *slog.Logger) *Server {
	s := &Server{
		Handlers: make(map[string]grafanaJSONServer.Handler),
		//reportsCache: reportsStore,
	}

	options := []grafanaJSONServer.Option{
		grafanaJSONServer.WithLogger(logger.With("component", "grafana-json-server")),
		grafanaJSONServer.WithRequestLogger(slog.LevelDebug, middleware.DefaultRequestLogFormatter),
		grafanaJSONServer.WithHTTPHandler(http.MethodGet, "/health", http.HandlerFunc(s.Health)),
		grafanaJSONServer.WithPrometheusQueryMetrics("sciensano", "", "sciensano"),
	}

	queryHandler := queryhandlers.SummaryHandler{ReportsStore: reportsStore}
	for _, m := range []grafanaJSONServer.Metric{
		newMetric("cases", sciensano.CasesValidSummaryModes().List()...),
		newMetric("hospitalisations", sciensano.HospitalisationsValidSummaryModes().List()...),
		newMetric("mortalities", sciensano.MortalitiesValidSummaryModes().List()...),
		newMetric("tests", sciensano.TestResultsValidSummaryModes().List()...),
	} {
		s.Handlers[m.Value] = queryHandler
		options = append(options, grafanaJSONServer.WithMetric(m, queryHandler, nil))

	}
	queryHandler = queryhandlers.SummaryHandler{ReportsStore: reportsStore, Accumulate: true}
	for _, m := range []grafanaJSONServer.Metric{
		newMetric("vaccinations-rate-Partial", sciensano.VaccinationsValidSummaryModes().List()...),
		newMetric("vaccinations-rate-Full", sciensano.VaccinationsValidSummaryModes().List()...),
	} {
		s.Handlers[m.Value] = queryHandler
		options = append(options, grafanaJSONServer.WithMetric(m, queryHandler, nil))

	}

	//{Fetch: s.vaccinationRate, Accumulate: true, Metric: newMetric("vaccinations-rate", byAgeGroup, byRegion)},

	/*
		for _, h := range []*Handler2{
			{Fetch: s.vaccinationFilteredRate, Accumulate: true, DoseType: sciensano.Partial, Metric: newMetric("vaccinations-rate-partial", byAgeGroup, byRegion)},
			{Fetch: s.vaccinationFilteredRate, Accumulate: true, DoseType: sciensano.Full, Metric: newMetric("vaccinations-rate-full", byAgeGroup, byRegion)},
		} {
			s.handlers[h.Metric.Value] = h
			options = append(options, grafanaJSONServer.WithMetric(h.Metric, h, nil))
		}
	*/
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
