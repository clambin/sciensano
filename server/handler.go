package server

import (
	"github.com/clambin/go-common/set"
	"github.com/clambin/sciensano/internal/sciensano"
)

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
