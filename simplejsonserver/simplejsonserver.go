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
func (server *Server) Initialize(ctx context.Context) {
	// set up handler lookup table
	server.initializeHandlers()
	// set up auto-refresh of demographics
	go server.Demographics.Run(ctx)
}

// Run runs the API handler server
func (server *Server) Run(port int) (err error) {
	server.Initialize(context.Background())
	r := server.GetRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	r.Path("/health").Handler(http.HandlerFunc(server.Health))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (server *Server) Health(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Handlers      int
		ReporterCache map[string]int
	}{
		Handlers:      len(server.Handlers),
		ReporterCache: server.Reporter.ReportCache.Stats(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
}

func (server *Server) initializeHandlers() {
	server.Handlers = map[string]simplejson.Handler{
		"cases":                     &cases.Handler{Reporter: server.Reporter, Scope: cases.ScopeAll},
		"cases-province":            &cases.Handler{Reporter: server.Reporter, Scope: cases.ScopeProvince},
		"cases-region":              &cases.Handler{Reporter: server.Reporter, Scope: cases.ScopeRegion},
		"cases-age":                 &cases.Handler{Reporter: server.Reporter, Scope: cases.ScopeAge},
		"hospitalisations":          &hospitalisations.Handler{Reporter: server.Reporter, Scope: hospitalisations.ScopeAll},
		"hospitalisations-region":   &hospitalisations.Handler{Reporter: server.Reporter, Scope: hospitalisations.ScopeRegion},
		"hospitalisations-province": &hospitalisations.Handler{Reporter: server.Reporter, Scope: hospitalisations.ScopeProvince},
		"mortality":                 &mortality.Handler{Reporter: server.Reporter, Scope: mortality.ScopeAll},
		"mortality-region":          &mortality.Handler{Reporter: server.Reporter, Scope: mortality.ScopeRegion},
		"mortality-age":             &mortality.Handler{Reporter: server.Reporter, Scope: mortality.ScopeAge},
		"tests":                     &testresults.Handler{Reporter: server.Reporter},
		"vaccinations":              &vaccinations.Handler{Reporter: server.Reporter},
		"vaccination-lag":           &vaccinations.LagHandler{Reporter: server.Reporter},
		"vacc-age-partial":          &vaccinations.GroupedHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypePartial, Accumulate: true},
		"vacc-age-full":             &vaccinations.GroupedHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeFull, Accumulate: true},
		"vacc-age-booster":          &vaccinations.GroupedHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeBooster},
		"vacc-region-partial":       &vaccinations.GroupedHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypePartial, Accumulate: true},
		"vacc-region-full":          &vaccinations.GroupedHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeFull, Accumulate: true},
		"vacc-region-booster":       &vaccinations.GroupedHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeBooster},
		"vacc-age-rate-partial":     &vaccinations.RateHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypePartial, Fetcher: server.Demographics},
		"vacc-age-rate-full":        &vaccinations.RateHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeFull, Fetcher: server.Demographics},
		"vacc-age-rate-booster":     &vaccinations.RateHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeAge, Type: vaccinations2.TypeBooster, Fetcher: server.Demographics},
		"vacc-region-rate-partial":  &vaccinations.RateHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypePartial, Fetcher: server.Demographics},
		"vacc-region-rate-full":     &vaccinations.RateHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeFull, Fetcher: server.Demographics},
		"vacc-region-rate-booster":  &vaccinations.RateHandler{Reporter: server.Reporter, Scope: vaccinations.ScopeRegion, Type: vaccinations2.TypeBooster, Fetcher: server.Demographics},
		"vacc-manufacturer":         &vaccinations.ManufacturerHandler{Reporter: server.Reporter},
		"vaccines":                  &vaccinesHandler.OverviewHandler{Reporter: server.Reporter},
		"vaccines-manufacturer":     &vaccinesHandler.ManufacturerHandler{Reporter: server.Reporter},
		"vaccines-stats":            &vaccinesHandler.StatsHandler{Reporter: server.Reporter},
		"vaccines-time":             &vaccinesHandler.DelayHandler{Reporter: server.Reporter},
	}
}
