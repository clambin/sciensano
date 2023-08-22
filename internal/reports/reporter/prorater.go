package reporter

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/v2/internal/population"
	"github.com/clambin/sciensano/v2/internal/population/bracket"
	"github.com/clambin/sciensano/v2/internal/reports/store"
	"github.com/clambin/sciensano/v2/internal/sciensano"
	"log/slog"
	"time"
)

type ProRater struct {
	Name     string
	Source   Publisher[sciensano.Vaccinations]
	PopStore PopulationFetcher
	Mode     sciensano.SummaryColumn
	DoseType sciensano.DoseType
	Store    *store.Store
	Logger   *slog.Logger
}

type PopulationFetcher interface {
	// GetByAgeBracket returns the demographics for the specified AgeBracket
	GetByAgeBracket(arguments bracket.Bracket) (count int)
	// GetByRegion returns the demographics grouped by region
	GetByRegion() (figures map[string]int)
	// WaitTillReady waits until the fetcher is ready or until the context is marked as done
	WaitTillReady(ctx context.Context) error
}

var _ PopulationFetcher = &population.Server{}

func (r *ProRater) Run(ctx context.Context) error {
	ch := make(chan sciensano.Vaccinations, 1)
	r.Source.Register(ch)
	defer func() {
		r.Source.Unregister(ch)
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-ch:
			r.createReport(data)
		}
	}
}

func (r *ProRater) createReport(vaccinations sciensano.Vaccinations) {
	t := tabulator.New()
	columnNames := set.Create[string]()

	// Filtering and then calling summary has a major performance impact.
	// This is basically the same code as Summary, but filters on the fly to avoid copying the large vaccinations slice.
	for _, vaccination := range vaccinations {
		if vaccination.Dose != r.DoseType && !(r.DoseType == sciensano.Full && vaccination.Dose == sciensano.SingleDose) {
			continue
		}

		columnName, err := vaccination.GetSummaryColumnName(r.Mode)
		if err != nil {
			r.Logger.Error("failed to generate report", "err", err)
			return
		}

		if columnName == "" {
			columnName = "(unknown)"
		}
		if !columnNames.Contains(columnName) {
			t.RegisterColumn(columnName)
			columnNames.Add(columnName)
		}

		t.Add(vaccination.TimeStamp.Time, columnName, float64(vaccination.Count))
	}

	t, err := proRate(t, r.Mode, r.PopStore)
	if err != nil {
		r.Logger.Error("failed to generate prorated report", "err", err)
		return
	}
	r.Store.Put(r.Name, t)
}

func proRate(summary *tabulator.Tabulator, mode sciensano.SummaryColumn, popStore PopulationFetcher) (*tabulator.Tabulator, error) {
	figures, err := getPopulationForGroup(mode, summary.GetColumns(), popStore)
	if err != nil {
		return nil, err
	}
	rated := summary.Copy()
	timestamps := rated.GetTimestamps()
	for _, column := range rated.GetColumns() {
		values, _ := rated.GetValues(column)
		for index, oldValue := range values {
			var newValue float64
			figure, ok := figures[column]
			if ok && figure != 0 {
				newValue = oldValue / float64(figure)
			}
			rated.Set(timestamps[index], column, newValue)
		}
	}
	return rated, nil
}

func getPopulationForGroup(mode sciensano.SummaryColumn, columns []string, popStore PopulationFetcher) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := popStore.WaitTillReady(ctx); err != nil {
		return nil, fmt.Errorf("population figures not ready: %w", err)
	}
	var figures map[string]int
	switch mode {
	case sciensano.ByRegion:
		figures = popStore.GetByRegion()
	case sciensano.ByAgeGroup:
		figures = make(map[string]int)
		for _, column := range columns {
			if column == "(unknown)" {
				continue
			}
			b, err := bracket.FromString(column)
			if err != nil {
				return nil, fmt.Errorf("invalid age bracket: '%s' : %w", column, err)
			}
			figures[column] = popStore.GetByAgeBracket(b)
		}
	}
	return figures, nil
}
