package vaccinations

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/sciensano"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 vaccinations
type Handler struct {
	Sciensano    sciensano.APIClient
	Demographics demographics.Demographics
	targetTable  grafanajson.TargetTable
}

// New creates a new Handler
func New(client sciensano.APIClient, demographics demographics.Demographics) (handler *Handler) {
	handler = &Handler{
		Sciensano:    client,
		Demographics: demographics,
	}

	handler.targetTable = grafanajson.TargetTable{
		"vaccinations":             {TableQueryFunc: handler.buildVaccinationTableResponse},
		"vacc-age-partial":         {TableQueryFunc: handler.buildGroupedVaccinationTableResponse},
		"vacc-age-full":            {TableQueryFunc: handler.buildGroupedVaccinationTableResponse},
		"vacc-age-booster":         {TableQueryFunc: handler.buildGroupedVaccinationTableResponse},
		"vacc-age-rate-partial":    {TableQueryFunc: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-age-rate-full":       {TableQueryFunc: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-age-rate-booster":    {TableQueryFunc: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-region-partial":      {TableQueryFunc: handler.buildGroupedVaccinationTableResponse},
		"vacc-region-full":         {TableQueryFunc: handler.buildGroupedVaccinationTableResponse},
		"vacc-region-booster":      {TableQueryFunc: handler.buildGroupedVaccinationTableResponse},
		"vacc-region-rate-partial": {TableQueryFunc: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-region-rate-full":    {TableQueryFunc: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-region-rate-booster": {TableQueryFunc: handler.buildGroupedVaccinationRateTableResponse},
		"vaccination-lag":          {TableQueryFunc: handler.buildVaccinationLagTableResponse},
	}

	return
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() grafanajson.Endpoints {
	return grafanajson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

// Search implements the grafana-json Search function. It returns all supported targets
func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

// TableQuery implements the grafana-json TableQuery function. It processes incoming TableQuery requests
func (handler *Handler) TableQuery(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	start := time.Now()
	response, err = handler.targetTable.RunTableQuery(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("unable to build table query response for target '%s': %s", target, err.Error())
	}
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Debug("TableQuery called")
	return
}

func (handler *Handler) buildVaccinationTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var vaccinationData *sciensano.Vaccinations
	vaccinationData, err = handler.Sciensano.GetVaccinations(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("failed to create grouped vaccination table response: %s", err.Error())
	}

	sciensano.AccumulateVaccinations(vaccinationData)

	timeStampColumn := make(grafanajson.TableQueryResponseTimeColumn, len(vaccinationData.Timestamps))
	partialColumn := make(grafanajson.TableQueryResponseNumberColumn, len(vaccinationData.Timestamps))
	fullColumn := make(grafanajson.TableQueryResponseNumberColumn, len(vaccinationData.Timestamps))
	boosterColumn := make(grafanajson.TableQueryResponseNumberColumn, len(vaccinationData.Timestamps))

	for index, timestamp := range vaccinationData.Timestamps {
		timeStampColumn[index] = timestamp
		partialColumn[index] = float64(vaccinationData.Groups[0].Values[index].Partial)
		fullColumn[index] = float64(vaccinationData.Groups[0].Values[index].Full + vaccinationData.Groups[0].Values[index].SingleDose)
		boosterColumn[index] = float64(vaccinationData.Groups[0].Values[index].Booster)

	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timeStampColumn},
			{Text: "partial", Data: partialColumn},
			{Text: "full", Data: fullColumn},
			{Text: "booster", Data: boosterColumn},
		},
	}
	return
}

func (handler *Handler) buildGroupedVaccinationTableResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var vaccinationData *sciensano.Vaccinations

	if strings.HasPrefix(target, "vacc-age-") {
		vaccinationData, err = handler.Sciensano.GetVaccinationsByAgeGroup(ctx, args.Range.To)
	} else if strings.HasPrefix(target, "vacc-region-") {
		vaccinationData, err = handler.Sciensano.GetVaccinationsByRegion(ctx, args.Range.To)
	} else {
		err = fmt.Errorf("invalid target: " + target)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create grouped vaccination table response: %s", err.Error())
	}

	sciensano.AccumulateVaccinations(vaccinationData)

	timeStampColumn := make(grafanajson.TableQueryResponseTimeColumn, len(vaccinationData.Timestamps))
	for index, timestamp := range vaccinationData.Timestamps {
		timeStampColumn[index] = timestamp
	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timeStampColumn},
		},
	}

	for _, series := range vaccinationData.Groups {
		dataColumn := make(grafanajson.TableQueryResponseNumberColumn, len(series.Values))

		for index, entry := range series.Values {
			var value int
			if strings.HasSuffix(target, "-partial") {
				value = entry.Partial
			} else if strings.HasSuffix(target, "-full") {
				value = entry.Full + entry.SingleDose
			} else if strings.HasSuffix(target, "-booster") {
				value = entry.Booster
			}
			dataColumn[index] = float64(value)
		}

		name := series.Name
		if name == "" {
			name = "(unknown)"
		}

		response.Columns = append(response.Columns, grafanajson.TableQueryResponseColumn{
			Text: name,
			Data: dataColumn,
		})
	}

	return
}

func (handler *Handler) buildGroupedVaccinationRateTableResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	response, err = handler.buildGroupedVaccinationTableResponse(ctx, target, args)

	if err != nil {
		return nil, fmt.Errorf("failed to get grouped vaccination rate figures: %s", err.Error())
	}

	response.Columns = filterUnknownColumns(response.Columns)

	if strings.HasPrefix(target, "vacc-age-rate-") {
		ageGroupFigures := handler.Demographics.GetAgeGroupFigures()
		prorateFigures(response, ageGroupFigures)
	} else if strings.HasPrefix(target, "vacc-region-rate-") {
		regionFigures := handler.Demographics.GetRegionFigures()
		_, ok := regionFigures["Ostbelgien"]
		if !ok {
			population, _ := regionFigures["Wallonia"]
			population -= 78000
			regionFigures["Wallonia"] = population
			regionFigures["Ostbelgien"] = 78000
		}
		prorateFigures(response, regionFigures)
	}
	return
}

func filterUnknownColumns(columns []grafanajson.TableQueryResponseColumn) []grafanajson.TableQueryResponseColumn {
	newColumns := make([]grafanajson.TableQueryResponseColumn, 0, len(columns))
	shouldReplace := false
	for _, column := range columns {
		if column.Text == "(unknown)" {
			shouldReplace = true
			continue
		}
		newColumns = append(newColumns, column)
	}
	if shouldReplace {
		return newColumns
	}
	return columns
}

func prorateFigures(result *grafanajson.TableQueryResponse, groups map[string]int) {
	newColumns := make([]grafanajson.TableQueryResponseColumn, 0, len(result.Columns))
	for _, column := range result.Columns {
		// TODO: perform this in a go routine and use WaitGroup to wait till done
		// set up a benchmark to check speed improvement
		if column.Text != "(empty)" {
			switch data := column.Data.(type) {
			case grafanajson.TableQueryResponseNumberColumn:
				figure, ok := groups[column.Text]
				for index, entry := range data {
					if ok && figure != 0 {
						data[index] = entry / float64(figure)
					} else {
						data[index] = 0
					}
				}
			}
			newColumns = append(newColumns, column)
		}
	}
	result.Columns = newColumns
}

func (handler *Handler) buildVaccinationLagTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var vaccinationsData *sciensano.Vaccinations

	vaccinationsData, err = handler.Sciensano.GetVaccinations(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("failed to determine vaccination lag: %s", err.Error())
	}

	sciensano.AccumulateVaccinations(vaccinationsData)
	timestamps, lag := buildLag(vaccinationsData)

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "lag", Data: lag},
	}

	return
}
