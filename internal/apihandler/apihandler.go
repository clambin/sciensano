package apihandler

import (
	"fmt"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"time"
)

// Handler implements a Grafana SimpleJson API Handler that gets BE covid stats
type Handler struct {
	Cache *cache.Cache

	lastDate map[string]time.Time
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
		"vaccination-lag",
	}
}

// Query the DB and return the requested targets
func (handler *Handler) Query(target string, _ *grafana_json.QueryRequest) (response *grafana_json.QueryResponse, err error) {
	err = fmt.Errorf("dataserie not implemented for '%s'", target)
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

	case "vaccination-lag":
		req := cache.VaccinationsRequest{
			EndTime:  request.Range.To,
			Response: make(chan []sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.Response
		vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
		response = buildVaccinationLagTableResponse(vaccineStats)

	default:
		err = fmt.Errorf("unknown target '%s'", target)
	}

	// log if there's a new update
	if err == nil && response != nil {
		handler.logUpdates(target, response)
	}

	return
}

func (handler *Handler) logUpdates(target string, response *grafana_json.QueryTableResponse) {
	if handler.lastDate == nil {
		handler.lastDate = make(map[string]time.Time)
	}

	key := "vaccinations"
	if target == "tests" {
		key = "tests"
	}

	latest := time.Time{}
	for _, column := range response.Columns {
		if column.Text == "timestamp" {
			timestamps := column.Data.(grafana_json.QueryTableResponseTimeColumn)
			if len(timestamps) > 0 {
				latest = timestamps[len(timestamps)-1]
				break
			}
		}
	}

	if entry, ok := handler.lastDate[key]; ok == false {
		handler.lastDate[key] = latest
	} else if latest.After(entry) {
		handler.lastDate[key] = latest
		log.WithFields(log.Fields{"target": key, "time": latest}).Info("new data found")
	}
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
	rows := len(vaccinations)
	timestamps := make(grafana_json.QueryTableResponseTimeColumn, rows)
	partial := make(grafana_json.QueryTableResponseNumberColumn, rows)
	full := make(grafana_json.QueryTableResponseNumberColumn, rows)

	for index, entry := range vaccinations {
		timestamps[index] = entry.Timestamp
		partial[index] = float64(entry.FirstDose)
		full[index] = float64(entry.SecondDose)
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
	// use goroutines to perform this in parallel
	results := make(map[string]chan []float64)
	for group := range vaccinations {
		results[group] = make(chan []float64)
		go func(groupName string, channel chan []float64) {
			channel <- getFilledVaccinations(timestamps, vaccinations[groupName], complete)
		}(group, results[group])
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
	for _, group := range groups {
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

func buildVaccinationLagTableResponse(vaccinations []sciensano.Vaccination) (response *grafana_json.QueryTableResponse) {
	vaccinationLag := buildLag(vaccinations)
	rows := len(vaccinationLag)

	timestamps := make(grafana_json.QueryTableResponseTimeColumn, rows)
	lag := make(grafana_json.QueryTableResponseNumberColumn, rows)

	for index, entry := range buildLag(vaccinations) {
		timestamps[index] = entry.Timestamp
		lag[index] = entry.Lag
	}

	response = new(grafana_json.QueryTableResponse)
	response.Columns = []grafana_json.QueryTableResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "lag", Data: lag},
	}

	return
}
