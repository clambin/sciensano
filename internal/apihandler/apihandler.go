package apihandler

import (
	"fmt"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/internal/demographics"
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Handler implements a Grafana SimpleJson API Handler that gets BE covid stats
type Handler struct {
	Cache    *cache.Cache
	Vaccines *vaccines.Handler

	lastDate map[string]time.Time
}

// Create a Handler
// TODO: this is a bit clunky. needs a better interface in grafana-json
func Create() (*Handler, error) {
	c := cache.New(15 * time.Minute)
	go c.Run()

	v := vaccines.Create()

	handler := Handler{
		Cache:    c,
		Vaccines: v,
	}

	return &handler, nil
}

// Endpoints tells the server which endpoints we have implemented
func (handler *Handler) Endpoints() grafana_json.Endpoints {
	return grafana_json.Endpoints{
		Search:      handler.Search,
		TableQuery:  handler.TableQuery,
		Annotations: handler.Annotations,
	}
}

// Search returns all supported targets
func (handler *Handler) Search() []string {
	return []string{
		"tests",
		"vaccinations",
		"vacc-age-partial", "vacc-age-full", "vacc-age-rate-partial", "vacc-age-rate-full",
		"vacc-region-partial", "vacc-region-full", "vacc-region-rate-partial", "vacc-region-rate-full",
		"vaccination-lag",
		"vaccines",
	}
}

func (handler *Handler) TableQuery(target string, args *grafana_json.TableQueryArgs) (response *grafana_json.TableQueryResponse, err error) {
	switch target {
	case "tests":
		req := cache.TestsRequest{
			EndTime:  args.Range.To,
			Response: make(chan []sciensano.Test),
		}
		handler.Cache.Tests <- req
		testStats := <-req.Response
		response = buildTestTableResponse(testStats)
	case "vaccinations":
		req := cache.VaccinationsRequest{
			EndTime:  args.Range.To,
			Response: make(chan []sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.Response
		vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
		response = buildVaccinationTableResponse(vaccineStats)

	case "vacc-age-partial", "vacc-age-full", "vacc-age-rate-partial", "vacc-age-rate-full":
		req := cache.VaccinationsRequest{
			EndTime:         args.Range.To,
			Filter:          "AgeGroup",
			GroupedResponse: make(chan map[string][]sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.GroupedResponse
		for ageGroup, data := range vaccineStats {
			vaccineStats[ageGroup] = sciensano.AccumulateVaccinations(data)
		}
		response = buildGroupedVaccinationTableResponse(vaccineStats, target)

		if strings.HasPrefix(target, "vacc-age-rate-") {
			prorateFigures(response, demographics.GetAgeGroupFigures())
		}

	case "vacc-region-partial", "vacc-region-full", "vacc-region-rate-partial", "vacc-region-rate-full":
		req := cache.VaccinationsRequest{
			EndTime:         args.Range.To,
			Filter:          "Region",
			GroupedResponse: make(chan map[string][]sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.GroupedResponse
		for region, data := range vaccineStats {
			vaccineStats[region] = sciensano.AccumulateVaccinations(data)
		}
		response = buildGroupedVaccinationTableResponse(vaccineStats, target)

		if strings.HasPrefix(target, "vacc-region-rate-") {
			prorateFigures(response, demographics.GetRegionFigures())
		}

	case "vaccination-lag":
		req := cache.VaccinationsRequest{
			EndTime:  args.Range.To,
			Response: make(chan []sciensano.Vaccination),
		}
		handler.Cache.Vaccinations <- req
		vaccineStats := <-req.Response
		vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
		response = buildVaccinationLagTableResponse(vaccineStats)

	case "vaccines":
		responseChannel := make(vaccines.ResponseChannel)
		handler.Vaccines.Request <- responseChannel
		batches := <-responseChannel
		batches = vaccines.AccumulateBatches(batches)
		response = buildVaccineTableResponse(batches)

	default:
		err = fmt.Errorf("unknown target '%s'", target)
	}

	// log if there's a new update
	if err == nil && response != nil {
		handler.logUpdates(target, response)
	}

	return
}

func (handler *Handler) logUpdates(target string, response *grafana_json.TableQueryResponse) {
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
			timestamps := column.Data.(grafana_json.TableQueryResponseTimeColumn)
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

func buildTestTableResponse(tests []sciensano.Test) (response *grafana_json.TableQueryResponse) {
	rows := len(tests)
	timestamps := make(grafana_json.TableQueryResponseTimeColumn, rows)
	allTests := make(grafana_json.TableQueryResponseNumberColumn, rows)
	positiveTests := make(grafana_json.TableQueryResponseNumberColumn, rows)
	positiveRate := make(grafana_json.TableQueryResponseNumberColumn, rows)

	for index, test := range tests {
		timestamps[index] = test.Timestamp
		allTests[index] = float64(test.Total)
		positiveTests[index] = float64(test.Positive)
		positiveRate[index] = float64(test.Positive) / float64(test.Total)
	}

	response = new(grafana_json.TableQueryResponse)
	response.Columns = []grafana_json.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "total", Data: allTests},
		{Text: "positive", Data: positiveTests},
		{Text: "rate", Data: positiveRate},
	}
	return
}

func buildVaccinationTableResponse(vaccinations []sciensano.Vaccination) (response *grafana_json.TableQueryResponse) {
	// TODO: pre-allocating size is more efficient
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

	return
}

func buildGroupedVaccinationTableResponse(vaccinations map[string][]sciensano.Vaccination, target string) (response *grafana_json.TableQueryResponse) {
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
	timestampColumn := make(grafana_json.TableQueryResponseTimeColumn, 0)
	dataColumns := make(map[string]grafana_json.TableQueryResponseNumberColumn, len(vaccinations))
	for _, group := range groups {
		dataColumns[group] = make(grafana_json.TableQueryResponseNumberColumn, 0)
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
	return
}

func buildVaccinationLagTableResponse(vaccinations []sciensano.Vaccination) (response *grafana_json.TableQueryResponse) {
	vaccinationLag := buildLag(vaccinations)
	rows := len(vaccinationLag)

	timestamps := make(grafana_json.TableQueryResponseTimeColumn, rows)
	lag := make(grafana_json.TableQueryResponseNumberColumn, rows)

	for index, entry := range buildLag(vaccinations) {
		timestamps[index] = entry.Timestamp
		lag[index] = entry.Lag
	}

	response = new(grafana_json.TableQueryResponse)
	response.Columns = []grafana_json.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "lag", Data: lag},
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

func buildVaccineTableResponse(batches []vaccines.Batch) (response *grafana_json.TableQueryResponse) {
	rows := len(batches)
	timestampColumn := make(grafana_json.TableQueryResponseTimeColumn, rows)
	batchColumn := make(grafana_json.TableQueryResponseNumberColumn, rows)

	for index, entry := range batches {
		timestampColumn[index] = time.Time(entry.Date)
		batchColumn[index] = float64(entry.Amount)
	}

	response = new(grafana_json.TableQueryResponse)
	response.Columns = []grafana_json.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestampColumn},
		{Text: "vaccines", Data: batchColumn},
	}
	return
}

func (handler *Handler) Annotations(name, query string, args *grafana_json.AnnotationRequestArgs) (annotations []grafana_json.Annotation, err error) {
	log.WithFields(log.Fields{
		"name":    name,
		"query":   query,
		"endTime": args.Range.To,
	}).Info("annotations")

	responseChannel := make(vaccines.ResponseChannel)
	handler.Vaccines.Request <- responseChannel

	batches := <-responseChannel
	for _, batch := range batches {
		annotations = append(annotations, grafana_json.Annotation{
			Time:  time.Time(batch.Date),
			Title: batch.Manufacturer,
			Text:  "Amount: " + strconv.FormatInt(batch.Amount, 10),
		})
	}
	return
}
