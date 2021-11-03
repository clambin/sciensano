package apihandler

import (
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	casesHandler "github.com/clambin/sciensano/apihandler/cases"
	mortalityHandler "github.com/clambin/sciensano/apihandler/mortality"
	covidTestsHandler "github.com/clambin/sciensano/apihandler/testresults"
	vaccinationsHandler "github.com/clambin/sciensano/apihandler/vaccinations"
	vaccinesHandler "github.com/clambin/sciensano/apihandler/vaccines"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	"net/http"
	"time"
)

// Handlers groups Grafana SimpleJson API handlers that retrieve Belgium COVID-19-related statistics
type Handlers struct {
	Sciensano    sciensano.APIClient
	Vaccines     vaccines.APIClient
	Demographics demographics.Demographics
	handlers     []grafanajson.Handler
}

// Create a Handlers object
func Create() *Handlers {
	handler := Handlers{
		Sciensano: sciensano.NewCachedClient(15 * time.Minute),
		Vaccines: &vaccines.Cache{
			APIClient: &vaccines.Client{HTTPClient: &http.Client{}},
			Retention: 24 * time.Hour,
		},
		Demographics: &demographics.Store{
			Retention:   24 * time.Hour,
			AgeBrackets: demographics.DefaultAgeBrackets,
		},
	}

	// force load of demographics data on startup
	go handler.Demographics.GetRegionFigures()

	handler.handlers = []grafanajson.Handler{
		covidTestsHandler.New(handler.Sciensano),
		vaccinationsHandler.New(handler.Sciensano, handler.Demographics),
		vaccinesHandler.New(handler.Sciensano, handler.Vaccines),
		casesHandler.New(handler.Sciensano),
		mortalityHandler.New(handler.Sciensano),
	}

	return &handler
}

// GetHandlers returns all configured handlers
func (handler *Handlers) GetHandlers() []grafanajson.Handler {
	return handler.handlers
}

// Run runs the API handler server
func (handler *Handlers) Run(port int) (err error) {
	server := grafanajson.Server{
		Handlers: handler.GetHandlers(),
	}
	r := server.GetRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
