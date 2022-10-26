package simplejsonserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/vaccinations"
	"github.com/clambin/sciensano/simplejsonserver/booster"
	vaccinations2 "github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
	"net/http"
)

// Server groups Grafana SimpleJson API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	simplejson.Server
	Reporter     *reporter.Client
	Demographics demographics.Fetcher
}

// Initialize starts background tasks to support Server
func (s *Server) Initialize(ctx context.Context) {
	// set up auto-refresh of demographics
	go s.Demographics.Run(ctx)
	// set up handler lookup table

	b := booster.Handler{Reporter: s.Reporter}

	s.Handlers = map[string]simplejson.Handler{
		"cases":                     &Handler{Fetcher: s.Reporter.Cases.Get},
		"cases-province":            &Handler{Fetcher: s.Reporter.Cases.GetByProvince},
		"cases-region":              &Handler{Fetcher: s.Reporter.Cases.GetByRegion},
		"cases-age":                 &Handler{Fetcher: s.Reporter.Cases.GetByAgeGroup},
		"hospitalisations":          &Handler{Fetcher: s.Reporter.Hospitalisations.Get},
		"hospitalisations-region":   &Handler{Fetcher: s.Reporter.Hospitalisations.GetByRegion},
		"hospitalisations-province": &Handler{Fetcher: s.Reporter.Hospitalisations.GetByProvince},
		"mortality":                 &Handler{Fetcher: s.Reporter.Mortality.Get},
		"mortality-region":          &Handler{Fetcher: s.Reporter.Mortality.GetByRegion},
		"mortality-age":             &Handler{Fetcher: s.Reporter.Mortality.GetByAgeGroup},
		"tests":                     &Handler{Fetcher: s.Reporter.TestResults.Get},
		"vaccinations":              &vaccinations2.Handler{Reporter: s.Reporter},
		"vacc-age-rate-partial":     &vaccinations2.RateHandler{Reporter: s.Reporter, Scope: vaccinations2.ByAge, Type: vaccinations.TypePartial, Fetcher: s.Demographics},
		"vacc-age-rate-full":        &vaccinations2.RateHandler{Reporter: s.Reporter, Scope: vaccinations2.ByAge, Type: vaccinations.TypeFull, Fetcher: s.Demographics},
		"vacc-region-rate-partial":  &vaccinations2.RateHandler{Reporter: s.Reporter, Scope: vaccinations2.ByRegion, Type: vaccinations.TypePartial, Fetcher: s.Demographics},
		"vacc-region-rate-full":     &vaccinations2.RateHandler{Reporter: s.Reporter, Scope: vaccinations2.ByRegion, Type: vaccinations.TypeFull, Fetcher: s.Demographics},
		"vacc-age-booster":          &vaccinations2.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations2.ByAge, Type: vaccinations.TypeBooster},
		"vacc-region-booster":       &vaccinations2.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations2.ByRegion, Type: vaccinations.TypeBooster},
		"vacc-manufacturer":         &Handler{Fetcher: s.Reporter.Vaccinations.GetByManufacturer, Accumulate: true},
		"boosters":                  &Handler{Fetcher: b.Fetch, Accumulate: true},
	}
}

// Run runs the API handler server
func (s *Server) Run(port int) (err error) {
	s.Initialize(context.Background())
	r := s.Server.GetRouter()
	r.HandleFunc("/health", s.Health)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Handlers      int
		ReporterCache map[string]int
	}{
		Handlers:      len(s.Handlers),
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
	return simplejson.Endpoints{
		Query: func(_ context.Context, req query.Request) (query.Response, error) {
			records, err := h.Fetcher()
			if err != nil {
				return nil, fmt.Errorf("fetch failed: %w", err)
			}
			if h.Accumulate {
				records = records.Accumulate()
			}
			return records.Filter(req.Args).CreateTableResponse(), nil
		},
	}
}
