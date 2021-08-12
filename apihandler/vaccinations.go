package apihandler

import (
	"context"
	"fmt"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"strings"
	"time"
)

func (handler *Handler) buildVaccinationTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse) {
	if vaccinations, err := handler.Sciensano.GetVaccinations(ctx, endTime); err == nil {
		vaccinations = sciensano.AccumulateVaccinations(vaccinations)

		rows := len(vaccinations)
		timestamps := make(grafanaJson.TableQueryResponseTimeColumn, rows)
		partial := make(grafanaJson.TableQueryResponseNumberColumn, rows)
		full := make(grafanaJson.TableQueryResponseNumberColumn, rows)

		for index, entry := range vaccinations {
			timestamps[index] = entry.Timestamp
			partial[index] = float64(entry.FirstDose)
			full[index] = float64(entry.SecondDose)
		}

		response = new(grafanaJson.TableQueryResponse)
		response.Columns = []grafanaJson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "partial", Data: partial},
			{Text: "full", Data: full},
		}
	}
	return
}

func (handler *Handler) buildGroupedVaccinationTableResponse(ctx context.Context, _, endTime time.Time, target string) (response *grafanaJson.TableQueryResponse) {
	var vaccinations map[string][]sciensano.Vaccination
	var err error

	if strings.HasPrefix(target, "vacc-age-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByAge(ctx, endTime)
	} else if strings.HasPrefix(target, "vacc-region-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByRegion(ctx, endTime)
	} else {
		err = fmt.Errorf("invalid target: " + target)
	}

	if err == nil {
		// grouped vaccinations are shown incrementally
		for ageGroup, data := range vaccinations {
			vaccinations[ageGroup] = sciensano.AccumulateVaccinations(data)
		}

		// sort group names so they always show up in the same order
		groups := getGroups(vaccinations)

		// get all timestamps across all groups & populate the timestamp column
		timestamps := getTimestamps(vaccinations)
		timestampCount := len(timestamps)

		// fill out each group, so all groups have all timestamps
		results := fillVaccinations(timestamps, vaccinations, strings.HasSuffix(target, "-full"))

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
	}
	return
}

func (handler *Handler) buildGroupedVaccinationRateTableResponse(ctx context.Context, beginTime, endTime time.Time, target string) (response *grafanaJson.TableQueryResponse) {
	response = handler.buildGroupedVaccinationTableResponse(ctx, beginTime, endTime, target)
	if response != nil {
		if strings.HasPrefix(target, "vacc-age-rate-") {
			prorateFigures(response, handler.demographics.GetAgeGroupFigures())
		} else if strings.HasPrefix(target, "vacc-region-rate-") {
			prorateFigures(response, handler.demographics.GetRegionFigures())
		}
	}
	return
}

func (handler *Handler) buildVaccinationLagTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse) {
	if vaccinations, err := handler.Sciensano.GetVaccinations(ctx, endTime); err == nil {
		vaccinations = sciensano.AccumulateVaccinations(vaccinations)

		timestamps, lag := buildLag(vaccinations)

		response = new(grafanaJson.TableQueryResponse)
		response.Columns = []grafanaJson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "lag", Data: lag},
		}
	}
	return
}

func prorateFigures(result *grafanaJson.TableQueryResponse, groups map[string]int) {
	newColumns := make([]grafanaJson.TableQueryResponseColumn, 0, len(result.Columns))
	for _, column := range result.Columns {
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
