package apihandler

import (
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	covidTestsHandler "github.com/clambin/sciensano/apihandler/covidtests"
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

// Create a Handlers
func Create() *Handlers {
	handler := Handlers{
		Sciensano: &sciensano.Client{
			Getter: &apiclient.Cache{
				Getter:    &apiclient.Client{HTTPClient: &http.Client{}},
				Retention: 15 * time.Minute,
			},
		},
		Vaccines: &vaccines.Cache{
			APIClient: &vaccines.Client{HTTPClient: &http.Client{}},
			Retention: time.Hour,
		},
		Demographics: &demographics.Store{
			Retention:   24 * time.Hour,
			AgeBrackets: demographics.DefaultAgeBrackets,
		},
	}

	handler.handlers = []grafanajson.Handler{
		covidTestsHandler.New(handler.Sciensano),
		vaccinationsHandler.New(handler.Sciensano, handler.Demographics),
		vaccinesHandler.New(handler.Sciensano, handler.Vaccines),
	}

	return &handler
}

func (handler *Handlers) GetHandlers() []grafanajson.Handler {
	return handler.handlers
}
