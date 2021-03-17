package apihandler

import (
	"errors"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/pkg/sciensano"
	"sort"
	"strings"
	"time"
)

// Handler implements a Grafana SimpleJson API Handler that gets BE covid stats
type Handler struct {
	Cache *cache.Cache
}

// Create a Handler
func Create() (*Handler, error) {
	c := cache.New(15 * time.Minute)
	go c.Run()
	return &Handler{
		Cache: c,
	}, nil
}

// Search returns all supported targets
func (handler *Handler) Search() []string {
	return []string{
		"tests",
		"vaccinations",
		"vacc-age-partial", "vacc-age-full",
		"vacc-region-partial", "vacc-region-full",
	}
}

// Query the DB and return the requested targets
func (handler *Handler) Query(target string, _ *grafana_json.QueryRequest) (response *grafana_json.QueryResponse, err error) {
	err = errors.New("dataserie not implemented for " + target)
	return
}

func (handler *Handler) QueryTable(target string, request *grafana_json.QueryRequest) (response *grafana_json.QueryTableResponse, err error) {
	switch target {
	case "tests":
		req := cache.TestsRequest{
			EndTime:  request.Range.To,
			Response: make(chan []sciensano.Test),
		}
		handler.Cache.Tests <- req
		testStats := <-req.Response
		response = buildTestTableResponse(testStats)
	case "vaccinations":
		req := cache.VaccinationsRequest{
			EndTime:  request.Range.To,
			Response: make(chan []sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.Response
		vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
		response = buildVaccinationTableResponse(vaccineStats)

	case "vacc-age-partial", "vacc-age-full":
		req := cache.VaccinationsRequest{
			EndTime:         request.Range.To,
			Filter:          "AgeGroup",
			GroupedResponse: make(chan map[string][]sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.GroupedResponse
		for ageGroup, data := range vaccineStats {
			vaccineStats[ageGroup] = sciensano.AccumulateVaccinations(data)
		}
		response = buildGroupedVaccinationTableResponse(vaccineStats, target)

	case "vacc-region-partial", "vacc-region-full":
		req := cache.VaccinationsRequest{
			EndTime:         request.Range.To,
			Filter:          "Region",
			GroupedResponse: make(chan map[string][]sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.GroupedResponse
		for region, data := range vaccineStats {
			vaccineStats[region] = sciensano.AccumulateVaccinations(data)
		}
		response = buildGroupedVaccinationTableResponse(vaccineStats, target)
	}

	return
}

func buildTestTableResponse(tests []sciensano.Test) (response *grafana_json.QueryTableResponse) {
	rows := len(tests)
	timestamps := make(grafana_json.QueryTableResponseTimeColumn, rows)
	allTests := make(grafana_json.QueryTableResponseNumberColumn, rows)
	positiveTests := make(grafana_json.QueryTableResponseNumberColumn, rows)
	positiveRate := make(grafana_json.QueryTableResponseNumberColumn, rows)

	for index, test := range tests {
		timestamps[index] = test.Timestamp
		allTests[index] = float64(test.Total)
		positiveTests[index] = float64(test.Positive)
		positiveRate[index] = float64(test.Positive) / float64(test.Total)
	}

	response = new(grafana_json.QueryTableResponse)
	response.Columns = []grafana_json.QueryTableResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "total", Data: allTests},
		{Text: "positive", Data: positiveTests},
		{Text: "rate", Data: positiveRate},
	}
	return
}

func buildVaccinationTableResponse(vaccinations []sciensano.Vaccination) (response *grafana_json.QueryTableResponse) {
	// TODO: pre-allocating size is more efficient
	timestamps := make(grafana_json.QueryTableResponseTimeColumn, 0)
	partial := make(grafana_json.QueryTableResponseNumberColumn, 0)
	full := make(grafana_json.QueryTableResponseNumberColumn, 0)

	for _, entry := range vaccinations {
		timestamps = append(timestamps, entry.Timestamp)
		partial = append(partial, float64(entry.FirstDose))
		full = append(full, float64(entry.SecondDose))
	}

	response = new(grafana_json.QueryTableResponse)
	response.Columns = []grafana_json.QueryTableResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "partial", Data: partial},
		{Text: "full", Data: full},
	}

	return
}

func buildGroupedVaccinationTableResponse(vaccinations map[string][]sciensano.Vaccination, target string) (response *grafana_json.QueryTableResponse) {
	// sort group names so they always show up in the same order
	groups := make([]string, len(vaccinations))
	index := 0
	for group := range vaccinations {
		groups[index] = group
		index++
	}
	sort.Strings(groups)

	// build the columns
	// TODO: pre-allocating size is more efficient
	timestampColumn := make(grafana_json.QueryTableResponseTimeColumn, 0)
	dataColumns := make(map[string]grafana_json.QueryTableResponseNumberColumn, len(vaccinations))
	for _, group := range groups {
		dataColumns[group] = make(grafana_json.QueryTableResponseNumberColumn, 0)
	}

	// get all timestamps across all groups & populate the timestamp column
	timestamps := getTimestamps(vaccinations)
	timestampColumn = append(timestampColumn, timestamps...)

	// do we want partial or complete vaccination figures?
	complete := false
	if strings.HasSuffix(target, "-full") {
		complete = true
	}

	// fill out each group, so all groups have all timestamps
	// use goroutines to performs this in parallel and get the results when we build the response
	results := make(map[string]chan []float64)
	for group := range vaccinations {
		results[group] = make(chan []float64)
		go func(groupName string) {
			results[groupName] <- getFilledVaccinations(timestamps, vaccinations[groupName], complete)
		}(group)
	}

	// populate the data columns
	for group := range vaccinations {
		data := <-results[group]
		dataColumns[group] = append(dataColumns[group], data...)
	}

	// build the response
	response = new(grafana_json.QueryTableResponse)
	response.Columns = []grafana_json.QueryTableResponseColumn{{
		Text: "timestamp",
		Data: timestampColumn,
	}}
	for group := range dataColumns {
		label := group
		if label == "" {
			label = "(empty)"
		}
		response.Columns = append(response.Columns, grafana_json.QueryTableResponseColumn{
			Text: label,
			Data: dataColumns[group],
		})
	}
	return
}

func getTimestamps(vaccinations map[string][]sciensano.Vaccination) (timestamps []time.Time) {
	// get unique timestamps
	uniqueTimestamps := make(map[time.Time]bool, len(vaccinations))
	for _, groupData := range vaccinations {
		for _, data := range groupData {
			uniqueTimestamps[data.Timestamp] = true
		}
	}
	timestamps = make([]time.Time, 0, len(uniqueTimestamps))
	for timestamp := range uniqueTimestamps {
		timestamps = append(timestamps, timestamp)
	}
	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i].Before(timestamps[j]) })
	return
}

// getFilledVaccinations has two main goals: 1) it returns either the partial or complete vaccination figures for a group and
// 2) it fills out the series with any missing timestamps so all columns cover the complete range of timestamps
func getFilledVaccinations(timestamps []time.Time, vaccinations []sciensano.Vaccination, complete bool) (filled []float64) {
	timestampCount := len(timestamps)
	vaccinationCount := len(vaccinations)

	var timestampIndex, vaccinationIndex int
	var lastVaccination sciensano.Vaccination

	for timestampIndex < timestampCount {
		for vaccinationIndex < vaccinationCount && timestamps[timestampIndex].Before(vaccinations[vaccinationIndex].Timestamp) {
			lastVaccination.Timestamp = timestamps[timestampIndex]
			filled = append(filled, float64(getVaccination(lastVaccination, complete)))
			timestampIndex++
		}
		if vaccinationIndex < vaccinationCount && timestamps[timestampIndex].Equal(vaccinations[vaccinationIndex].Timestamp) {
			lastVaccination = vaccinations[vaccinationIndex]
			filled = append(filled, float64(getVaccination(lastVaccination, complete)))
			vaccinationIndex++
			timestampIndex++
		} else if vaccinationIndex == vaccinationIndex {
			for ; timestampIndex < timestampCount; timestampIndex++ {
				lastVaccination.Timestamp = timestamps[timestampIndex]
				filled = append(filled, float64(getVaccination(lastVaccination, complete)))
			}
		}
	}
	return
}

func getVaccination(vaccination sciensano.Vaccination, complete bool) (value int) {
	if complete == false {
		value = vaccination.FirstDose
	} else {
		value = vaccination.SecondDose
	}
	return
}
