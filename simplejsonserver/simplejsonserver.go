package simplejsonserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/reporter"
	vaccinations2 "github.com/clambin/sciensano/reporter/vaccinations"
	"github.com/clambin/sciensano/simplejsonserver/cases"
	"github.com/clambin/sciensano/simplejsonserver/hospitalisations"
	"github.com/clambin/sciensano/simplejsonserver/mortality"
	"github.com/clambin/sciensano/simplejsonserver/testresults"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	vaccinesHandler "github.com/clambin/sciensano/simplejsonserver/vaccines"
	"github.com/clambin/simplejson/v3"
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
	s.Handlers = map[string]simplejson.Handler{
		"cases":                     &cases.Handler{Reporter: s.Reporter, Scope: cases.ScopeAll},
		"cases-province":            &cases.Handler{Reporter: s.Reporter, Scope: cases.ScopeProvince},
		"cases-region":              &cases.Handler{Reporter: s.Reporter, Scope: cases.ScopeRegion},
		"cases-age":                 &cases.Handler{Reporter: s.Reporter, Scope: cases.ScopeAge},
		"hospitalisations":          &hospitalisations.Handler{Reporter: s.Reporter, Scope: hospitalisations.ScopeAll},
		"hospitalisations-region":   &hospitalisations.Handler{Reporter: s.Reporter, Scope: hospitalisations.ScopeRegion},
		"hospitalisations-province": &hospitalisations.Handler{Reporter: s.Reporter, Scope: hospitalisations.ScopeProvince},
		"mortality":                 &mortality.Handler{Reporter: s.Reporter, Scope: mortality.ScopeAll},
		"mortality-region":          &mortality.Handler{Reporter: s.Reporter, Scope: mortality.ScopeRegion},
		"mortality-age":             &mortality.Handler{Reporter: s.Reporter, Scope: mortality.ScopeAge},
		"tests":                     &testresults.Handler{Reporter: s.Reporter},
		"vaccinations":              &vaccinations.Handler{Reporter: s.Reporter},
		"vaccination-lag":           &vaccinations.LagHandler{Reporter: s.Reporter},
		"vacc-age-partial":          &vaccinations.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypePartial, Accumulate: true},
		"vacc-age-full":             &vaccinations.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeFull, Accumulate: true},
		"vacc-age-booster":          &vaccinations.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeBooster},
		"vacc-region-partial":       &vaccinations.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypePartial, Accumulate: true},
		"vacc-region-full":          &vaccinations.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeFull, Accumulate: true},
		"vacc-region-booster":       &vaccinations.GroupedHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeBooster},
		"vacc-age-rate-partial":     &vaccinations.RateHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypePartial, Fetcher: s.Demographics},
		"vacc-age-rate-full":        &vaccinations.RateHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeFull, Fetcher: s.Demographics},
		"vacc-age-rate-booster":     &vaccinations.RateHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeBooster, Fetcher: s.Demographics},
		"vacc-region-rate-partial":  &vaccinations.RateHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypePartial, Fetcher: s.Demographics},
		"vacc-region-rate-full":     &vaccinations.RateHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeFull, Fetcher: s.Demographics},
		"vacc-region-rate-booster":  &vaccinations.RateHandler{Reporter: s.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeBooster, Fetcher: s.Demographics},
		"vacc-manufacturer":         &vaccinations.ManufacturerHandler{Reporter: s.Reporter},
		"vaccines":                  &vaccinesHandler.OverviewHandler{Reporter: s.Reporter},
		"vaccines-manufacturer":     &vaccinesHandler.ManufacturerHandler{Reporter: s.Reporter},
		"vaccines-stats":            &vaccinesHandler.StatsHandler{Reporter: s.Reporter},
		"vaccines-time":             &vaccinesHandler.DelayHandler{Reporter: s.Reporter},
	}
}

// Run runs the API handler server
func (s *Server) Run(port int) (err error) {
	s.Initialize(context.Background())
	r := s.Server.GetRouter()
	r.Path("/health").Handler(http.HandlerFunc(s.Health))
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
