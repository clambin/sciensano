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
	// GetForAgeBracket returns the population for the specified AgeBracket
	GetForAgeBracket(bracket bracket.Bracket) (count int)
	// GetForRegion returns the population region
	GetForRegion(region string) (count int)
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
	t, err := filterVaccinations(vaccinations, r.Mode, r.DoseType)
	if err != nil {
		r.Logger.Error("failed to generate report", "err", err)
		return
	}
	t, err = proRate(t, r.Mode, r.PopStore)
	if err != nil {
		r.Logger.Error("failed to generate prorated report", "err", err)
		return
	}
	r.Store.Put(r.Name, t)
}

func filterVaccinations(vaccinations sciensano.Vaccinations, mode sciensano.SummaryColumn, doseType sciensano.DoseType) (*tabulator.Tabulator, error) {
	t := tabulator.New()
	columnNames := set.New[string]()

	// Filtering and then calling summary has a major performance impact.
	// This is basically the same code as Summary, but filters on the fly to avoid copying the large vaccinations slice.
	for i := range vaccinations {
		if vaccinations[i].Dose != doseType && !(doseType == sciensano.Full && vaccinations[i].Dose == sciensano.SingleDose) {
			continue
		}

		columnName, err := vaccinations[i].GetSummaryColumnName(mode)
		if err != nil {
			return nil, err
		}

		if columnName == "" {
			columnName = "(unknown)"
		}
		if !columnNames.Contains(columnName) {
			t.RegisterColumn(columnName)
			columnNames.Add(columnName)
		}

		t.Add(vaccinations[i].TimeStamp.Time, columnName, float64(vaccinations[i].Count))
	}
	return t, nil
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

	figures := make(map[string]int)

	for _, column := range columns {
		if column == "(unknown)" {
			continue
		}
		var pop int
		switch mode {
		case sciensano.ByRegion:
			pop = popStore.GetForRegion(column)
		case sciensano.ByAgeGroup:
			b, err := bracket.FromString(column)
			if err != nil {
				return nil, fmt.Errorf("invalid age bracket: '%s' : %w", column, err)
			}
			pop = popStore.GetForAgeBracket(b)
		default:
		}
		figures[column] = pop
	}
	return figures, nil
}
