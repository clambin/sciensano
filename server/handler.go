package server

import (
	"github.com/clambin/go-common/set"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
)

func buildHandlers(s ReportsStore) []*SummaryHandler {
	summaryHandlers := []struct {
		name           string
		summaryColumns set.Set[sciensano.SummaryColumn]
	}{
		{name: "cases", summaryColumns: sciensano.CasesValidSummaryModes()},
		{name: "hospitalisations", summaryColumns: sciensano.HospitalisationsValidSummaryModes()},
		{name: "mortalities", summaryColumns: sciensano.MortalitiesValidSummaryModes()},
		{name: "tests", summaryColumns: sciensano.TestResultsValidSummaryModes()},
		{name: "vaccinations", summaryColumns: sciensano.VaccinationsValidSummaryModes()},
	}

	var handlers []*SummaryHandler
	for _, h := range summaryHandlers {
		processor := summaryMetricProcessor{summaryColumns: h.summaryColumns}
		handlers = append(handlers, &SummaryHandler{
			ReportsStore:    s,
			Metric:          grafanaJSONServer.Metric{Label: h.name, Value: h.name, Payloads: processor.makeMetricPayload()},
			metricProcessor: processor,
		})
	}

	return handlers
}
