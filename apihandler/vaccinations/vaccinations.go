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
	var vaccinations []sciensano.Vaccination
	vaccinations, err = handler.Sciensano.GetVaccinations(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("unable to get vaccination data: %s", err.Error())
	}

	vaccinations = sciensano.AccumulateVaccinations(vaccinations)

	rows := len(vaccinations)
	timestamps := make(grafanajson.TableQueryResponseTimeColumn, rows)
	partial := make(grafanajson.TableQueryResponseNumberColumn, rows)
	full := make(grafanajson.TableQueryResponseNumberColumn, rows)
	booster := make(grafanajson.TableQueryResponseNumberColumn, rows)

	for index, entry := range vaccinations {
		timestamps[index] = entry.Timestamp
		partial[index] = float64(entry.Partial)
		full[index] = float64(entry.Full + entry.SingleDose)
		booster[index] = float64(entry.Booster)
	}

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "partial", Data: partial},
		{Text: "full", Data: full},
		{Text: "booster", Data: booster},
	}
	return
}

func (handler *Handler) buildGroupedVaccinationTableResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var vaccinations map[string][]sciensano.Vaccination

	if strings.HasPrefix(target, "vacc-age-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByAge(ctx, args.Range.To)
	} else if strings.HasPrefix(target, "vacc-region-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByRegion(ctx, args.Range.To)
	} else {
		err = fmt.Errorf("invalid target: " + target)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create grouped vaccination table response: %s", err.Error())
	}

	// grouped vaccinations are shown incrementally
	for group, data := range vaccinations {
		vaccinations[group] = sciensano.AccumulateVaccinations(data)
	}

	// sort group names so they always show up in the same order
	groups := getGroups(vaccinations)

	// get all timestamps across all groups & populate the timestamp column
	timestamps := getTimestamps(vaccinations)
	timestampCount := len(timestamps)

	// fill out each group, so all groups have all timestamps
	results := fillVaccinations(timestamps, vaccinations, getGroupType(target))

	// build & populate the timestamp columns
	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, timestampCount)
	timestampColumn = append(timestampColumn, timestamps...)

	// build & populate the data columns
	dataColumns := make(map[string]grafanajson.TableQueryResponseNumberColumn, len(groups))
	for _, group := range groups {
		dataColumns[group] = make(grafanajson.TableQueryResponseNumberColumn, 0, timestampCount)
		data := <-results[group]
		dataColumns[group] = append(dataColumns[group], data...)
	}

	// build the response
	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{{
		Text: "timestamp",
		Data: timestampColumn,
	}}
	for _, group := range groups {
		label := group
		if label == "" {
			label = "(empty)"
		}
		response.Columns = append(response.Columns, grafanajson.TableQueryResponseColumn{
			Text: label,
			Data: dataColumns[group],
		})
	}
	return
}

func (handler *Handler) buildGroupedVaccinationRateTableResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	response, err = handler.buildGroupedVaccinationTableResponse(ctx, target, args)

	if err != nil {
		return nil, fmt.Errorf("failed to get grouped vaccination rate figures: %s", err.Error())
	}

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
	var vaccinations []sciensano.Vaccination
	vaccinations, err = handler.Sciensano.GetVaccinations(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("failed to determine vaccination lag: %s", err.Error())
	}

	vaccinations = sciensano.AccumulateVaccinations(vaccinations)
	timestamps, lag := buildLag(vaccinations)

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "lag", Data: lag},
	}

	return
}
