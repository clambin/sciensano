package apihandler

import (
	"fmt"
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/demographics"
	"github.com/clambin/sciensano/pkg/sciensano"
	"strings"
	"time"
)

func (handler *Handler) buildVaccinationTableResponse(endTime time.Time, _ string) (response *grafana_json.TableQueryResponse) {
	if vaccinations, err := handler.Sciensano.GetVaccinations(endTime); err == nil {
		vaccinations = sciensano.AccumulateVaccinations(vaccinations)

		rows := len(vaccinations)
		timestamps := make(grafana_json.TableQueryResponseTimeColumn, rows)
		partial := make(grafana_json.TableQueryResponseNumberColumn, rows)
		full := make(grafana_json.TableQueryResponseNumberColumn, rows)

		for index, entry := range vaccinations {
			timestamps[index] = entry.Timestamp
			partial[index] = float64(entry.FirstDose)
			full[index] = float64(entry.SecondDose)
		}

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "partial", Data: partial},
			{Text: "full", Data: full},
		}
	}
	return
}

func (handler *Handler) buildGroupedVaccinationTableResponse(endTime time.Time, target string) (response *grafana_json.TableQueryResponse) {
	var vaccinations map[string][]sciensano.Vaccination
	var err error

	if strings.HasPrefix(target, "vacc-age-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByAge(endTime)
	} else if strings.HasPrefix(target, "vacc-region-") {
		vaccinations, err = handler.Sciensano.GetVaccinationsByRegion(endTime)
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
		timestampColumn := make(grafana_json.TableQueryResponseTimeColumn, 0, timestampCount)
		timestampColumn = append(timestampColumn, timestamps...)

		// build & populate the data columns
		dataColumns := make(map[string]grafana_json.TableQueryResponseNumberColumn, len(groups))
		for _, group := range groups {
			dataColumns[group] = make(grafana_json.TableQueryResponseNumberColumn, 0, timestampCount)
			data := <-results[group]
			dataColumns[group] = append(dataColumns[group], data...)
		}

		// build the response
		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{{
			Text: "timestamp",
			Data: timestampColumn,
		}}
		for _, group := range groups {
			label := group
			if label == "" {
				label = "(empty)"
			}
			response.Columns = append(response.Columns, grafana_json.TableQueryResponseColumn{
				Text: label,
				Data: dataColumns[group],
			})
		}
	}
	return
}

func (handler *Handler) buildGroupedVaccinationRateTableResponse(endTime time.Time, target string) (response *grafana_json.TableQueryResponse) {
	response = handler.buildGroupedVaccinationTableResponse(endTime, target)
	if response != nil && strings.HasPrefix(target, "vacc-age-rate-") {
		prorateFigures(response, demographics.GetAgeGroupFigures())
	}
	return
}

func (handler *Handler) buildVaccinationLagTableResponse(endTime time.Time, _ string) (response *grafana_json.TableQueryResponse) {
	if vaccinations, err := handler.Sciensano.GetVaccinations(endTime); err == nil {
		vaccinations = sciensano.AccumulateVaccinations(vaccinations)

		timestamps, lag := buildLag(vaccinations)

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestamps},
			{Text: "lag", Data: lag},
		}
	}
	return
}

func prorateFigures(result *grafana_json.TableQueryResponse, groups map[string]int) {
	newColumns := make([]grafana_json.TableQueryResponseColumn, 0, len(result.Columns))
	for _, column := range result.Columns {
		if column.Text != "(empty)" {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseNumberColumn:
				if figure, ok := groups[column.Text]; ok && figure != 0 {
					for index, entry := range data {
						data[index] = entry / float64(figure)
					}
				}
			}
			newColumns = append(newColumns, column)
		}
	}
	result.Columns = newColumns
}
