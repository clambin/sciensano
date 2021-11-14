package vaccines

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler/response"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 cases
type Handler struct {
	Reporter    reporter.Reporter
	targetTable grafanajson.TargetTable
}

// New creates a new Handler
func New(reporter reporter.Reporter) (handler *Handler) {
	handler = &Handler{
		Reporter: reporter,
	}

	handler.targetTable = grafanajson.TargetTable{
		"vaccines":       {TableQueryFunc: handler.buildVaccineTableResponse},
		"vaccines-stats": {TableQueryFunc: handler.buildVaccineStatsTableResponse},
		"vaccines-time":  {TableQueryFunc: handler.buildVaccineTimeTableResponse},
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
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Info("TableQuery called")
	return
}

func (handler *Handler) buildVaccineTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (output *grafanajson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccines(ctx)

	if err == nil {
		batches.Accumulate()
		output = response.GenerateTableQueryResponse(batches, args)
	}

	return
}

func (handler *Handler) buildVaccineStatsTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccines(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	batches.Accumulate()
	batches.ApplyRange(args.Range.From, args.Range.To)

	var vaccinations *datasets.Dataset
	vaccinations, err = handler.Reporter.GetVaccinations(ctx)

	if err != nil {
		return
	}

	summarizeVaccinations(vaccinations)
	vaccinations.Accumulate()
	vaccinations.ApplyRange(args.Range.From, args.Range.To)

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn(vaccinations.Timestamps)},
		{Text: "vaccinations", Data: grafanajson.TableQueryResponseNumberColumn(vaccinations.Groups[0].Values)},
		{Text: "reserve", Data: grafanajson.TableQueryResponseNumberColumn(calculateVaccineReserve(vaccinations, batches))},
	}
	return
}

func summarizeVaccinations(vaccinationsData *datasets.Dataset) {
	for index := range vaccinationsData.Timestamps {
		total := 0.0
		for _, group := range vaccinationsData.Groups {
			total += group.Values[index]
		}
		vaccinationsData.Groups[0].Values[index] = total
	}
	vaccinationsData.Groups = vaccinationsData.Groups[:1]
}

func calculateVaccineReserve(vaccinationsData *datasets.Dataset, batches *datasets.Dataset) (reserve []float64) {
	batchIndex := 0
	lastBatch := 0.0

	for index, timestamp := range vaccinationsData.Timestamps {
		// find the last time we received vaccines
		for batchIndex < len(batches.Timestamps) &&
			!batches.Timestamps[batchIndex].After(timestamp) {
			// how many vaccines have we received so far?
			lastBatch = batches.Groups[0].Values[batchIndex]
			batchIndex++
		}

		// add it to the list
		reserve = append(reserve, lastBatch-vaccinationsData.Groups[0].Values[index])
	}

	return
}

func (handler *Handler) buildVaccineTimeTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccines(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	batches.Accumulate()
	batches.ApplyRange(args.Range.From, args.Range.To)

	var vaccinations *datasets.Dataset
	vaccinations, err = handler.Reporter.GetVaccinations(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccination data: %s", err.Error())
	}

	vaccinations.Accumulate()
	vaccinations.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn, timeColumn := calculateVaccineDelay(vaccinations, batches)

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestampColumn},
		{Text: "time", Data: timeColumn},
	}
	return
}

func calculateVaccineDelay(vaccinationsData *datasets.Dataset, batches *datasets.Dataset) (timestamps grafanajson.TableQueryResponseTimeColumn, delays grafanajson.TableQueryResponseNumberColumn) {
	batchIndex := 0
	batchCount := len(batches.Timestamps)

	for index, timestamp := range vaccinationsData.Timestamps {
		// how many vaccines did we need to perform this many vaccinations?
		vaccinesNeeded := vaccinationsData.Groups[0].Values[index] + vaccinationsData.Groups[1].Values[index]

		// find when we reached that number of vaccines
		for batchIndex < batchCount && batches.Groups[0].Values[batchIndex] < vaccinesNeeded {
			batchIndex++
		}

		// we depleted the *previous* batch. report the time difference between now and when we received that batch
		if batchIndex > 0 {
			timestamps = append(timestamps, timestamp)
			delays = append(delays, timestamp.Sub(batches.Timestamps[batchIndex-1]).Hours()/24)
		}
	}
	return
}
