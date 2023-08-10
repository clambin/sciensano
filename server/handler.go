package server

import (
	"github.com/clambin/go-common/set"
	"github.com/clambin/sciensano/internal/sciensano"
)

func buildHandlers(s ReportsStore) []SummaryHandler {
	summaryHandlers := []struct {
		name           string
		summaryColumns set.Set[sciensano.SummaryColumn]
		accumulate     bool
	}{
		{name: "cases", summaryColumns: sciensano.CasesValidSummaryModes()},
		{name: "hospitalisations", summaryColumns: sciensano.HospitalisationsValidSummaryModes()},
		{name: "mortalities", summaryColumns: sciensano.MortalitiesValidSummaryModes()},
		{name: "tests", summaryColumns: sciensano.TestResultsValidSummaryModes()},
		{name: "vaccinations", summaryColumns: sciensano.VaccinationsValidSummaryModes(), accumulate: true},
	}

	var handlers []SummaryHandler
	for _, h := range summaryHandlers {
		handlers = append(handlers, newSummaryHandler(h.name, h.summaryColumns, s))
	}

	return handlers
}

func buildDoseTypeHandlers(s ReportsStore) []SummaryByDoseTypeHandler {
	doseTypeHandlers := []struct {
		name           string
		summaryColumns set.Set[sciensano.SummaryColumn]
		doseTypes      set.Set[sciensano.DoseType]
		accumulate     bool
	}{
		{
			name:           "vaccinations-rate",
			summaryColumns: sciensano.VaccinationsValidSummaryModes(),
			doseTypes:      set.Create(sciensano.Partial, sciensano.Full),
		},
	}

	var handlers []SummaryByDoseTypeHandler
	for _, h := range doseTypeHandlers {
		handlers = append(handlers, newSummaryByDoseTypeHandler(h.name, h.summaryColumns, h.doseTypes, s))
	}

	return handlers
}
