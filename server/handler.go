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
		{name: "vaccinations-rate-Partial", summaryColumns: sciensano.VaccinationsValidSummaryModes()},
		{name: "vaccinations-rate-Full", summaryColumns: sciensano.VaccinationsValidSummaryModes()},
	}

	var handlers []SummaryHandler
	for _, h := range summaryHandlers {
		handlers = append(handlers, newSummaryHandler(h.name, h.summaryColumns, s))
	}

	return handlers
}
