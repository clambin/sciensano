package apihandler

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/vaccines"
	casesHandler "github.com/clambin/sciensano/apihandler/cases"
	hospitalisationsHandler "github.com/clambin/sciensano/apihandler/hospitalisations"
	mortalityHandler "github.com/clambin/sciensano/apihandler/mortality"
	covidTestsHandler "github.com/clambin/sciensano/apihandler/testresults"
	vaccinationsHandler "github.com/clambin/sciensano/apihandler/vaccinations"
	vaccinesHandler "github.com/clambin/sciensano/apihandler/vaccines"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter"
	"net/http"
	"time"
)

// Server groups Grafana SimpleJson API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	Reporter     *reporter.Client
	Vaccines     vaccines.Getter
	Demographics demographics.Demographics
	handlers     []grafanajson.Handler
}

const refreshInterval = 1 * time.Hour

// NewServer a Server object
func NewServer() (server *Server) {
	server = &Server{
		Reporter: reporter.NewCachedClient(refreshInterval),
		Vaccines: &vaccines.Client{
			HTTPClient: &http.Client{},
			Cache:      measurement.Cache{Retention: 24 * time.Hour},
		},
		Demographics: &demographics.Store{
			Retention:   24 * time.Hour,
			AgeBrackets: demographics.DefaultAgeBrackets,
		},
	}

	// force load of demographics data on startup
	go server.Demographics.GetRegionFigures()
	// set up auto-refresh of reports
	go server.Reporter.Sciensano.AutoRefresh(context.Background())

	server.handlers = []grafanajson.Handler{
		covidTestsHandler.New(server.Reporter),
		vaccinationsHandler.New(server.Reporter, server.Demographics),
		vaccinesHandler.New(server.Reporter),
		casesHandler.New(server.Reporter),
		mortalityHandler.New(server.Reporter),
		hospitalisationsHandler.New(server.Reporter),
	}

	return server
}

// GetHandlers returns all configured handlers
func (server *Server) GetHandlers() []grafanajson.Handler {
	return server.handlers
}

// Run runs the API handler server
func (server *Server) Run(port int) (err error) {
	s := grafanajson.Server{
		Handlers: server.GetHandlers(),
	}
	r := s.GetRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
