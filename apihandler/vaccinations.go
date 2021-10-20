package apihandler

import (
	"context"
	"fmt"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"strings"
	"time"
)

func (handler *Handler) buildVaccinationTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse, err error) {
	var vaccinations []sciensano.Vaccination
	vaccinations, err = handler.Sciensano.GetVaccinations(ctx, endTime)

	if err != nil {
		return nil, fmt.Errorf("unable to get vaccination data: %s", err.Error())
	}

	vaccinations = sciensano.AccumulateVaccinations(vaccinations)

	rows := len(vaccinations)
	timestamps := make(grafanaJson.TableQueryResponseTimeColumn, rows)
	partial := make(grafanaJson.TableQueryResponseNumberColumn, rows)
	full := make(grafanaJson.TableQueryResponseNumberColumn, rows)
	booster := make(grafanaJson.TableQueryResponseNumberColumn, rows)

	for index, entry := range vaccinations {
		timestamps[index] = entry.Timestamp
		partial[index] = float64(entry.Partial)
		full[index] = float64(entry.Full + entry.SingleDose)
		booster[index] = float64(entry.Booster)
	}

	response = new(grafanaJson.TableQueryResponse)
	response.Columns = []grafanaJson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "partial", Data: partial},
		{Text: "full", Data: full},
		{Text: "booster", Data: booster},
	}
	return
}

func (handler *Handler) buildGroupedVaccinationTableResponse(ctx context.Context, _, endTime time.Time, target string) (response *grafanaJson.TableQueryResponse, err error) {
	var vaccinations map[string][]sciensano.Vaccination

	if strings.HasPrefix(target, "vacc-age-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByAge(ctx, endTime)
	} else if strings.HasPrefix(target, "vacc-region-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByRegion(ctx, endTime)
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
	timestampColumn := make(grafanaJson.TableQueryResponseTimeColumn, 0, timestampCount)
	timestampColumn = append(timestampColumn, timestamps...)

	// build & populate the data columns
	dataColumns := make(map[string]grafanaJson.TableQueryResponseNumberColumn, len(groups))
	for _, group := range groups {
		dataColumns[group] = make(grafanaJson.TableQueryResponseNumberColumn, 0, timestampCount)
		data := <-results[group]
		dataColumns[group] = append(dataColumns[group], data...)
	}

	// build the response
	response = new(grafanaJson.TableQueryResponse)
	response.Columns = []grafanaJson.TableQueryResponseColumn{{
		Text: "timestamp",
		Data: timestampColumn,
	}}
	for _, group := range groups {
		label := group
		if label == "" {
			label = "(empty)"
		}
		response.Columns = append(response.Columns, grafanaJson.TableQueryResponseColumn{
			Text: label,
			Data: dataColumns[group],
		})
	}
	return
}

const (
	groupPartial = iota
	groupFull
	groupBooster
)

func getGroupType(target string) (groupType int) {
	if strings.HasSuffix(target, "-full") {
		groupType = groupFull
	} else if strings.HasSuffix(target, "-booster") {
		groupType = groupBooster
	} else {
		groupType = groupPartial
	}
	return
}

func (handler *Handler) buildGroupedVaccinationRateTableResponse(ctx context.Context, beginTime, endTime time.Time, target string) (response *grafanaJson.TableQueryResponse, err error) {
	response, err = handler.buildGroupedVaccinationTableResponse(ctx, beginTime, endTime, target)

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

func (handler *Handler) buildVaccinationLagTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse, err error) {
	var vaccinations []sciensano.Vaccination
	vaccinations, err = handler.Sciensano.GetVaccinations(ctx, endTime)

	if err != nil {
		return nil, fmt.Errorf("failed to determine vaccination lag: %s", err.Error())
	}

	vaccinations = sciensano.AccumulateVaccinations(vaccinations)
	timestamps, lag := buildLag(vaccinations)

	response = new(grafanaJson.TableQueryResponse)
	response.Columns = []grafanaJson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "lag", Data: lag},
	}

	return
}

func prorateFigures(result *grafanaJson.TableQueryResponse, groups map[string]int) {
	newColumns := make([]grafanaJson.TableQueryResponseColumn, 0, len(result.Columns))
	for _, column := range result.Columns {
		// TODO: perform this in a go routine and use WaitGroup to wait till done
		// set up a benchmark to check speed improvement
		if column.Text != "(empty)" {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseNumberColumn:
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
