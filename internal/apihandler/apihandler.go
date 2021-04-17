package apihandler

import (
	"errors"
	"fmt"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/demographics"
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Handler implements a Grafana SimpleJson API Handler that gets BE covid stats
type Handler struct {
	Sciensano sciensano.API
	Vaccines  *vaccines.Server

	lastDate map[string]time.Time
}

// Create a Handler
func Create() (*Handler, error) {
	handler := Handler{
		Sciensano: &sciensano.Client{
			HTTPClient:    &http.Client{Timeout: 20 * time.Second},
			CacheDuration: 15 * time.Minute,
		},
		Vaccines: vaccines.New(),
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
		"vaccines-stats",
	}
}

func (handler *Handler) TableQuery(target string, args *grafana_json.TableQueryArgs) (response *grafana_json.TableQueryResponse, err error) {
	switch target {
	case "tests":
		response = handler.buildTestTableResponse(args.Range.To)

	case "vaccinations":
		response = handler.buildVaccinationTableResponse(args.Range.To)

	case "vacc-age-partial", "vacc-age-full", "vacc-age-rate-partial", "vacc-age-rate-full":
		response = handler.buildGroupedVaccinationTableResponse(args.Range.To, target)
		if strings.HasPrefix(target, "vacc-age-rate-") {
			prorateFigures(response, demographics.GetAgeGroupFigures())
		}

	case "vacc-region-partial", "vacc-region-full", "vacc-region-rate-partial", "vacc-region-rate-full":
		response = handler.buildGroupedVaccinationTableResponse(args.Range.To, target)
		if strings.HasPrefix(target, "vacc-region-rate-") {
			prorateFigures(response, demographics.GetRegionFigures())
		}

	case "vaccination-lag":
		response = handler.buildVaccinationLagTableResponse(args.Range.To)

	case "vaccines":
		response = handler.buildVaccineTableResponse(args.Range.To)

	case "vaccines-stats":
		response = handler.buildVaccineStatsTableResponse(args.Range.To)

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

func (handler *Handler) buildTestTableResponse(endTime time.Time) (response *grafana_json.TableQueryResponse) {
	if tests, err := handler.Sciensano.GetTests(endTime); err == nil {

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
	}
	return
}

func (handler *Handler) buildVaccinationTableResponse(endTime time.Time) (response *grafana_json.TableQueryResponse) {
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
		err = errors.New("invalid target: " + target)
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

func (handler *Handler) buildVaccinationLagTableResponse(endTime time.Time) (response *grafana_json.TableQueryResponse) {
	if vaccinations, err := handler.Sciensano.GetVaccinations(endTime); err == nil {
		vaccinations = sciensano.AccumulateVaccinations(vaccinations)

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

func (handler *Handler) buildVaccineTableResponse(_ time.Time) (response *grafana_json.TableQueryResponse) {
	if batches, err := handler.Vaccines.GetBatches(); err == nil {
		batches = vaccines.AccumulateBatches(batches)

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
	}
	return
}

func (handler *Handler) buildVaccineStatsTableResponse(endTime time.Time) (response *grafana_json.TableQueryResponse) {
	var batches []vaccines.Batch
	var vaccinations []sciensano.Vaccination
	var err error
	if batches, err = handler.Vaccines.GetBatches(); err == nil {
		batches = vaccines.AccumulateBatches(batches)
		if vaccinations, err = handler.Sciensano.GetVaccinations(endTime); err == nil {
			vaccinations = sciensano.AccumulateVaccinations(vaccinations)
		}

		rows := len(vaccinations)
		timestampColumn := make(grafana_json.TableQueryResponseTimeColumn, 0, rows)
		reserveColumn := make(grafana_json.TableQueryResponseNumberColumn, 0, rows)
		delayColumn := make(grafana_json.TableQueryResponseNumberColumn, 0, rows)

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			for _, entry := range vaccinations {
				timestampColumn = append(timestampColumn, entry.Timestamp)
			}
			wg.Done()
		}()

		go func() {
			for _, value := range calculateVaccineReserve(vaccinations, batches) {
				reserveColumn = append(reserveColumn, value)
			}
			wg.Done()
		}()

		go func() {
			for _, value := range calculateVaccineDelay(vaccinations, batches) {
				delayColumn = append(delayColumn, value)
			}
			wg.Done()
		}()

		wg.Wait()

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "reserve", Data: reserveColumn},
			{Text: "delay", Data: delayColumn},
		}
	}
	return
}

func (handler *Handler) Annotations(name, query string, args *grafana_json.AnnotationRequestArgs) (annotations []grafana_json.Annotation, err error) {
	log.WithFields(log.Fields{
		"name":    name,
		"query":   query,
		"endTime": args.Range.To,
	}).Info("annotations")

	var batches []vaccines.Batch
	if batches, err = handler.Vaccines.GetBatches(); err == nil {
		for _, batch := range batches {
			annotations = append(annotations, grafana_json.Annotation{
				Time: time.Time(batch.Date),
				// Title: batch.Manufacturer,
				Text: "Amount: " + strconv.Itoa(batch.Amount),
			})
		}
	}
	return
}
