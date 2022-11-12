package simplejsonserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/httpserver"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/vaccinations"
	vaccinations2 "github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
	"net/http"
)

// Server groups Grafana SimpleJson API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	server       *simplejson.Server
	handlers     map[string]simplejson.Handler
	Reporter     *reporter.Client
	Demographics demographics.Fetcher
}

func New(port int, r *reporter.Client, f demographics.Fetcher) (s *Server, err error) {
	s = &Server{
		Reporter:     r,
		Demographics: f,
		handlers: map[string]simplejson.Handler{
			"cases":                     &Handler{Fetcher: r.Cases.Get},
			"cases-province":            &Handler{Fetcher: r.Cases.GetByProvince},
			"cases-region":              &Handler{Fetcher: r.Cases.GetByRegion},
			"cases-age":                 &Handler{Fetcher: r.Cases.GetByAgeGroup},
			"hospitalisations":          &Handler{Fetcher: r.Hospitalisations.Get},
			"hospitalisations-region":   &Handler{Fetcher: r.Hospitalisations.GetByRegion},
			"hospitalisations-province": &Handler{Fetcher: r.Hospitalisations.GetByProvince},
			"mortality":                 &Handler{Fetcher: r.Mortality.Get},
			"mortality-region":          &Handler{Fetcher: r.Mortality.GetByRegion},
			"mortality-age":             &Handler{Fetcher: r.Mortality.GetByAgeGroup},
			"tests":                     &Handler{Fetcher: r.TestResults.Get},
			"vaccinations":              &vaccinations2.Handler{Reporter: r},
			"vacc-age-rate-partial":     &vaccinations2.RateHandler{Reporter: r, Scope: vaccinations2.ByAge, Type: vaccinations.TypePartial, Fetcher: f},
			"vacc-age-rate-full":        &vaccinations2.RateHandler{Reporter: r, Scope: vaccinations2.ByAge, Type: vaccinations.TypeFull, Fetcher: f},
			"vacc-region-rate-partial":  &vaccinations2.RateHandler{Reporter: r, Scope: vaccinations2.ByRegion, Type: vaccinations.TypePartial, Fetcher: f},
			"vacc-region-rate-full":     &vaccinations2.RateHandler{Reporter: r, Scope: vaccinations2.ByRegion, Type: vaccinations.TypeFull, Fetcher: f},
			"vacc-age-booster":          &vaccinations2.GroupedHandler{Reporter: r, Scope: vaccinations2.ByAge, Type: vaccinations.TypeBooster},
			"vacc-region-booster":       &vaccinations2.GroupedHandler{Reporter: r, Scope: vaccinations2.ByRegion, Type: vaccinations.TypeBooster},
			"vacc-manufacturer":         &Handler{Fetcher: r.Vaccinations.GetByManufacturer, Accumulate: true},
			"boosters":                  &vaccinations2.BoosterHandler{Reporter: r},
		},
	}

	s.server, err = simplejson.New("sciensano", s.handlers,
		httpserver.WithPort{Port: port},
		httpserver.WithHandlers{Handlers: []httpserver.Handler{{
			Path:    "/health",
			Handler: http.HandlerFunc(s.Health),
			Methods: []string{http.MethodGet},
		}}},
	)
	return s, err
}

// Run runs the API handler server
func (s *Server) Run(ctx context.Context) (err error) {
	go s.Demographics.Run(ctx)
	return s.server.Run()
}

func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Handlers      int
		ReporterCache map[string]int
	}{
		Handlers:      len(s.handlers),
		ReporterCache: s.Reporter.ReportCache.Stats(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
}

type Handler struct {
	Fetcher    func() (*data.Table, error)
	Accumulate bool
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (h *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: h.query}
}

func (h *Handler) query(_ context.Context, req query.Request) (query.Response, error) {
	records, err := h.Fetcher()
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	if h.Accumulate {
		records = records.Accumulate()
	}
	return records.Filter(req.Args).CreateTableResponse(), nil
}
