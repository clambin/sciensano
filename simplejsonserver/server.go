package simplejsonserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/go-common/httpserver/middleware"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/cache"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/bracket"
	"github.com/clambin/sciensano/simplejsonserver/reports"
	"github.com/clambin/simplejson/v6"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

// Server groups Grafana SimpleJson API handlers that retrieve Belgium COVID-19-related statistics
type Server struct {
	server       *simplejson.Server
	handlers     map[string]simplejson.Handler
	apiCache     *cache.SciensanoCache
	reportsCache *reports.Cache
	Demographics demographics.Fetcher
}

var _ prometheus.Collector = &Server{}

func New(f demographics.Fetcher) (s *Server, err error) {
	s = &Server{
		apiCache:     cache.NewSciensanoCache(""),
		reportsCache: reports.NewCache(15 * time.Minute),
		Demographics: f,
	}
	s.handlers = map[string]simplejson.Handler{
		"cases":                     Handler{Fetch: s.cases, Mode: sciensano.Total},
		"cases-province":            Handler{Fetch: s.cases, Mode: sciensano.ByProvince},
		"cases-region":              Handler{Fetch: s.cases, Mode: sciensano.ByRegion},
		"cases-age":                 Handler{Fetch: s.cases, Mode: sciensano.ByAgeGroup},
		"hospitalisations":          Handler{Fetch: s.hospitalisations, Mode: sciensano.Total},
		"hospitalisations-region":   Handler{Fetch: s.hospitalisations, Mode: sciensano.ByRegion},
		"hospitalisations-province": Handler{Fetch: s.hospitalisations, Mode: sciensano.ByProvince},
		"mortality":                 Handler{Fetch: s.mortalities, Mode: sciensano.Total},
		"mortality-region":          Handler{Fetch: s.mortalities, Mode: sciensano.ByRegion},
		"mortality-age":             Handler{Fetch: s.mortalities, Mode: sciensano.ByAgeGroup},
		"tests":                     Handler{Fetch: s.testResults, Mode: sciensano.Total},
		"vaccinations":              Handler{Fetch: s.vaccinations, Mode: sciensano.Total, Accumulate: true},
		"vacc-age":                  Handler{Fetch: s.vaccinations, Mode: sciensano.ByAgeGroup, Accumulate: false},
		"vacc-region":               Handler{Fetch: s.vaccinations, Mode: sciensano.ByRegion, Accumulate: false},
		"vacc-manufacturer":         Handler{Fetch: s.vaccinations, Mode: sciensano.ByManufacturer, Accumulate: true},
		"vacc-region-rate":          Handler{Fetch: s.vaccinationRate, Mode: sciensano.ByRegion, Accumulate: true},
		"vacc-age-rate":             Handler{Fetch: s.vaccinationRate, Mode: sciensano.ByAgeGroup, Accumulate: true},
		"vacc-age-rate-partial":     Handler2{Fetch: s.vaccinationFilteredRate, Mode: sciensano.ByAgeGroup, DoseType: sciensano.Partial, Accumulate: true},
		"vacc-age-rate-full":        Handler2{Fetch: s.vaccinationFilteredRate, Mode: sciensano.ByAgeGroup, DoseType: sciensano.Full, Accumulate: true},
		"vacc-region-rate-partial":  Handler2{Fetch: s.vaccinationFilteredRate, Mode: sciensano.ByRegion, DoseType: sciensano.Partial, Accumulate: true},
		"vacc-region-rate-full":     Handler2{Fetch: s.vaccinationFilteredRate, Mode: sciensano.ByRegion, DoseType: sciensano.Full, Accumulate: true},
	}

	r := simplejson.New(s.handlers,
		simplejson.WithQueryMetrics{Name: "sciensano"},
		simplejson.WithHTTPMetrics{Option: middleware.PrometheusMetricsOptions{Application: "sciensano", MetricsType: middleware.Histogram}},
	)
	r.Get("/health", s.Health)
	s.server = r

	return s, err
}

// Serve runs the API handler server
func (s *Server) Serve(ctx context.Context, port int) (err error) {
	go s.Demographics.Run(ctx)
	go s.apiCache.AutoRefresh(ctx, time.Hour)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.server)
}

func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Handlers      int
		ReporterCache map[string]int
	}{
		Handlers:      len(s.handlers),
		ReporterCache: s.reportsCache.Stats(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
}

func (s *Server) cases(ctx context.Context, mode sciensano.SummaryColumn) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("cases-"+mode.String(), func() (*tabulator.Tabulator, error) {
		cases := s.apiCache.Cases.Get(ctx)
		return cases.Summarize(mode)
	})
}

func (s *Server) hospitalisations(ctx context.Context, mode sciensano.SummaryColumn) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("hospitalisations-"+mode.String(), func() (*tabulator.Tabulator, error) {
		hospitalisations := s.apiCache.Hospitalisations.Get(ctx)
		if mode == sciensano.Total {
			return hospitalisations.Categorize(), nil
		}
		return hospitalisations.Summarize(mode)
	})
}
func (s *Server) mortalities(ctx context.Context, mode sciensano.SummaryColumn) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("mortalities-"+mode.String(), func() (*tabulator.Tabulator, error) {
		mortalities := s.apiCache.Mortalities.Get(ctx)
		return mortalities.Summarize(mode)
	})
}

func (s *Server) testResults(ctx context.Context, mode sciensano.SummaryColumn) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("testResults-"+mode.String(), func() (*tabulator.Tabulator, error) {
		testResults := s.apiCache.TestResults.Get(ctx)
		if mode == sciensano.Total {
			return testResults.Categorize(), nil
		}
		return testResults.Summarize(mode)
	})
}

func (s *Server) vaccinations(ctx context.Context, mode sciensano.SummaryColumn) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("vaccinations-"+mode.String(), func() (*tabulator.Tabulator, error) {
		vaccinations := s.apiCache.Vaccinations.Get(ctx)
		if mode == sciensano.Total {
			return vaccinations.Categorize(), nil
		}
		return vaccinations.Summarize(mode)
	})
}

func (s *Server) vaccinationRate(ctx context.Context, mode sciensano.SummaryColumn) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("vaccinations-rate-"+mode.String(), func() (*tabulator.Tabulator, error) {
		vaccinations := s.apiCache.Vaccinations.Get(ctx)
		if mode == sciensano.Total {
			return nil, fmt.Errorf("rate not supported for Total mode")
		}
		v, err := vaccinations.Summarize(mode)
		if err != nil {
			return nil, err
		}

		var figures map[string]int
		switch mode {
		case sciensano.ByRegion:
			figures = s.Demographics.GetByRegion()
		case sciensano.ByAgeGroup:
			figures = make(map[string]int)
			for _, column := range v.GetColumns() {
				var b bracket.Bracket
				b, err = bracket.FromString(column)
				if err != nil {
					return nil, fmt.Errorf("invalid age bracket: '%s' : %w", column, err)
				}
				figures[column] = s.Demographics.GetByAgeBracket(b)
			}
		}
		return prorateFigures(v, figures), nil
	})
}

func (s *Server) vaccinationFilteredRate(ctx context.Context, mode sciensano.SummaryColumn, doseType sciensano.DoseType) (*tabulator.Tabulator, error) {
	return s.reportsCache.MaybeGenerate("vaccinations-filtered-rate-"+mode.String()+"-"+doseType.String(), func() (*tabulator.Tabulator, error) {
		vaccinations := s.apiCache.Vaccinations.Get(ctx)
		if mode == sciensano.Total {
			return vaccinations.Categorize(), nil
		}

		vaccinations = filterVaccinations(vaccinations, doseType)
		report, err := vaccinations.Summarize(mode)
		if err != nil {
			return nil, err
		}
		var figures map[string]int
		switch mode {
		case sciensano.ByRegion:
			figures = s.Demographics.GetByRegion()
		case sciensano.ByAgeGroup:
			figures = make(map[string]int)
			for _, column := range report.GetColumns() {
				var b bracket.Bracket
				b, err = bracket.FromString(column)
				if err != nil {
					return nil, fmt.Errorf("invalid age bracket: '%s' : %w", column, err)
				}
				figures[column] = s.Demographics.GetByAgeBracket(b)
			}
		}
		return prorateFigures(report, figures), nil
	})
}

func prorateFigures(d *tabulator.Tabulator, groups map[string]int) *tabulator.Tabulator {
	// old: Benchmark_Vaccinations-16    	  276972	      3710 ns/op
	// new: Benchmark_Vaccinations-16    	  283011	      3598 ns/op
	timestamps := d.GetTimestamps()
	for _, column := range d.GetColumns() {
		values, _ := d.GetValues(column)
		for index, oldValue := range values {
			var newValue float64
			figure, ok := groups[column]
			if ok {
				newValue = oldValue / float64(figure)
			}
			d.Set(timestamps[index], column, newValue)
		}
	}
	return d
}

func filterVaccinations(vaccinations sciensano.Vaccinations, doseType sciensano.DoseType) sciensano.Vaccinations {
	filtered := make(sciensano.Vaccinations, len(vaccinations))
	var index int
	for _, vaccination := range vaccinations {
		if vaccination.Dose == doseType || (doseType == sciensano.Full && vaccination.Dose == sciensano.SingleDose) {
			filtered[index] = vaccination
			index++
		}
	}
	filtered = filtered[:index]
	return filtered
}

func (s *Server) Describe(descs chan<- *prometheus.Desc) {
	s.apiCache.Describe(descs)
	s.server.Describe(descs)
}

func (s *Server) Collect(metrics chan<- prometheus.Metric) {
	s.apiCache.Collect(metrics)
	s.server.Collect(metrics)
}
