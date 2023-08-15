package reporter

import (
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/sciensano/internal/reports/datasource"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/sciensano"
	"log/slog"
)

type datasourceType int

const (
	casesDatasource datasourceType = iota
	hospitalisationsDatasource
	mortalitiesDatasource
	testResultsDatasource
	vaccinationsDatasource
)

func NewSciensanoReporters(datasources *datasource.SciensanoSources, store *store.Store, popStore Fetcher, logger *slog.Logger) []taskmanager.Task {
	summarizers := []struct {
		dsType   datasourceType
		basename string
		modes    []sciensano.SummaryColumn
	}{
		{dsType: casesDatasource, basename: "cases", modes: []sciensano.SummaryColumn{sciensano.Total, sciensano.ByProvince, sciensano.ByRegion, sciensano.ByAgeGroup}},
		{dsType: hospitalisationsDatasource, basename: "hospitalisations", modes: []sciensano.SummaryColumn{sciensano.Total, sciensano.ByProvince, sciensano.ByRegion, sciensano.ByCategory}},
		{dsType: mortalitiesDatasource, basename: "mortalities", modes: []sciensano.SummaryColumn{sciensano.Total, sciensano.ByRegion, sciensano.ByAgeGroup}},
		{dsType: testResultsDatasource, basename: "tests", modes: []sciensano.SummaryColumn{sciensano.Total, sciensano.ByCategory}},
		{dsType: vaccinationsDatasource, basename: "vaccinations", modes: []sciensano.SummaryColumn{sciensano.Total, sciensano.ByRegion, sciensano.ByAgeGroup, sciensano.ByManufacturer, sciensano.ByVaccinationType}},
	}

	var reporters []taskmanager.Task
	for _, option := range summarizers {
		for _, mode := range option.modes {
			fullName := option.basename + "-" + mode.String()
			l := logger.With(slog.String("reporter", fullName))

			var task taskmanager.Task
			switch option.dsType {
			case casesDatasource:
				task = &Summary[sciensano.Cases]{Name: fullName, Source: &datasources.Cases, Mode: mode, Store: store, Logger: l}
			case hospitalisationsDatasource:
				task = &Summary[sciensano.Hospitalisations]{Name: fullName, Source: &datasources.Hospitalisations, Mode: mode, Store: store, Logger: l}
			case mortalitiesDatasource:
				task = &Summary[sciensano.Mortalities]{Name: fullName, Source: &datasources.Mortalities, Mode: mode, Store: store, Logger: l}
			case testResultsDatasource:
				task = &Summary[sciensano.TestResults]{Name: fullName, Source: &datasources.TestResults, Mode: mode, Store: store, Logger: l}
			case vaccinationsDatasource:
				task = &Summary[sciensano.Vaccinations]{Name: fullName, Source: &datasources.Vaccinations, Mode: mode, Store: store, Logger: l}
			default:
				panic("invalid mode")
			}
			reporters = append(reporters, task)
		}
	}

	raters := []struct {
		dsType   datasourceType
		basename string
		modes    []sciensano.SummaryColumn
		doseType sciensano.DoseType
	}{
		{dsType: vaccinationsDatasource, basename: "vaccination-rate", modes: []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup}, doseType: sciensano.Partial},
		{dsType: vaccinationsDatasource, basename: "vaccination-rate", modes: []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup}, doseType: sciensano.Full},
	}

	for _, rater := range raters {
		for _, mode := range rater.modes {
			fullName := rater.basename + "-" + rater.doseType.String() + "-" + mode.String()
			l := logger.With("reporter", fullName)

			reporters = append(reporters, &ProRater{
				Name:     fullName,
				Source:   &datasources.Vaccinations,
				PopStore: popStore,
				Mode:     mode,
				DoseType: rater.doseType,
				Store:    store,
				Logger:   l,
			})
		}

	}

	return reporters
}
