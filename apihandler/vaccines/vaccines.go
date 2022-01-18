package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/apihandler/response"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 cases
type Handler struct {
	Reporter    reporter.Reporter
	targetTable simplejson.TargetTable
}

// New creates a new Handler
func New(reporter reporter.Reporter) (handler *Handler) {
	handler = &Handler{
		Reporter: reporter,
	}

	handler.targetTable = simplejson.TargetTable{
		"vaccines":              {TableQueryFunc: handler.buildVaccineTableResponse},
		"vaccines-manufacturer": {TableQueryFunc: handler.buildVaccineByManufacturerTableResponse},
		"vaccines-stats":        {TableQueryFunc: handler.buildVaccineStatsTableResponse},
		"vaccines-time":         {TableQueryFunc: handler.buildVaccineTimeTableResponse},
	}

	return
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

// Search implements the grafana-json Search function. It returns all supported targets
func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

// TableQuery implements the grafana-json TableQuery function. It processes incoming TableQuery requests
func (handler *Handler) TableQuery(ctx context.Context, target string, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	start := time.Now()
	response, err = handler.targetTable.RunTableQuery(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("unable to build table query response for target '%s': %s", target, err.Error())
	}
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Info("TableQuery called")
	return
}

func (handler *Handler) buildVaccineTableResponse(_ context.Context, _ string, args *simplejson.TableQueryArgs) (output *simplejson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccines()
	if err == nil {
		batches.Accumulate()
		output = response.GenerateTableQueryResponse(batches, args)
	}
	return
}

func (handler *Handler) buildVaccineByManufacturerTableResponse(_ context.Context, _ string, args *simplejson.TableQueryArgs) (output *simplejson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccinesByManufacturer()
	if err == nil {
		batches.Accumulate()
		output = response.GenerateTableQueryResponse(batches, args)
	}
	return
}

func (handler *Handler) buildVaccineStatsTableResponse(_ context.Context, _ string, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccines()

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	batches.Accumulate()

	var vaccinations *datasets.Dataset
	vaccinations, err = handler.Reporter.GetVaccinations()

	if err != nil {
		return
	}

	summarizeVaccinations(vaccinations)
	vaccinations.Accumulate()

	vaccinations.Groups = append(vaccinations.Groups, datasets.DatasetGroup{
		Name:   "reserve",
		Values: calculateVaccineReserve(vaccinations, batches),
	})

	batches.ApplyRange(args.Range.From, args.Range.To)
	vaccinations.ApplyRange(args.Range.From, args.Range.To)

	response = new(simplejson.TableQueryResponse)
	response.Columns = []simplejson.TableQueryResponseColumn{
		{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn(vaccinations.Timestamps)},
		{Text: "vaccinations", Data: simplejson.TableQueryResponseNumberColumn(vaccinations.Groups[0].Values)},
		{Text: "reserve", Data: simplejson.TableQueryResponseNumberColumn(vaccinations.Groups[1].Values)},
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

func (handler *Handler) buildVaccineTimeTableResponse(_ context.Context, _ string, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	var batches *datasets.Dataset
	batches, err = handler.Reporter.GetVaccines()

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	batches.Accumulate()
	batches.ApplyRange(args.Range.From, args.Range.To)

	var vaccinations *datasets.Dataset
	vaccinations, err = handler.Reporter.GetVaccinations()

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccination data: %s", err.Error())
	}

	vaccinations.Accumulate()
	vaccinations.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn, timeColumn := calculateVaccineDelay(vaccinations, batches)

	response = new(simplejson.TableQueryResponse)
	response.Columns = []simplejson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestampColumn},
		{Text: "time", Data: timeColumn},
	}
	return
}

func calculateVaccineDelay(vaccinationsData *datasets.Dataset, batches *datasets.Dataset) (timestamps simplejson.TableQueryResponseTimeColumn, delays simplejson.TableQueryResponseNumberColumn) {
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
