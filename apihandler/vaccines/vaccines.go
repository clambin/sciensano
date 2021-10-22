package vaccines

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Handler struct {
	Sciensano   sciensano.APIClient
	Vaccines    vaccines.APIClient
	targetTable grafanajson.TargetTable
}

func New(sciensanoClient sciensano.APIClient, vaccinesClient vaccines.APIClient) (handler *Handler) {
	handler = &Handler{
		Sciensano: sciensanoClient,
		Vaccines:  vaccinesClient,
	}

	handler.targetTable = grafanajson.TargetTable{
		"vaccines":       {TableQueryFunc: handler.buildVaccineTableResponse},
		"vaccines-stats": {TableQueryFunc: handler.buildVaccineStatsTableResponse},
		"vaccines-time":  {TableQueryFunc: handler.buildVaccineTimeTableResponse},
	}

	return
}

func (handler *Handler) Endpoints() grafanajson.Endpoints {
	return grafanajson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

func (handler *Handler) TableQuery(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	start := time.Now()
	response, err = handler.targetTable.RunTableQuery(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("unable to build table query response for target '%s': %s", target, err.Error())
	}
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Debug("TableQuery called")
	return
}

func (handler *Handler) buildVaccineTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var batches []*vaccines.Batch
	batches, err = handler.Vaccines.GetBatches(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	batches = vaccines.AccumulateBatches(batches)

	rows := len(batches)
	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, rows)
	batchColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, rows)

	for _, entry := range batches {
		if entry.Date.Time.After(args.Range.To) {
			continue
		}
		timestampColumn = append(timestampColumn, entry.Date.Time)
		batchColumn = append(batchColumn, float64(entry.Amount))
	}

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestampColumn},
		{Text: "vaccines", Data: batchColumn},
	}
	return
}

func (handler *Handler) buildVaccineStatsTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var batches []*vaccines.Batch
	batches, err = handler.Vaccines.GetBatches(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	var vaccinations []sciensano.Vaccination
	vaccinations, err = handler.Sciensano.GetVaccinations(ctx, args.Range.To)

	if err != nil {
		return
	}

	batches = vaccines.AccumulateBatches(batches)
	vaccinations = sciensano.AccumulateVaccinations(vaccinations)

	rows := len(vaccinations)
	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, rows)
	vaccinationsColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, rows)
	reserveColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, rows)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for _, entry := range vaccinations {
			timestampColumn = append(timestampColumn, entry.Timestamp)
			vaccinationsColumn = append(vaccinationsColumn, float64(entry.Total()))
		}
		wg.Done()
	}()

	go func() {
		for _, value := range calculateVaccineReserve(vaccinations, batches) {
			reserveColumn = append(reserveColumn, value)
		}
		wg.Done()
	}()

	wg.Wait()

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestampColumn},
		{Text: "vaccinations", Data: vaccinationsColumn},
		{Text: "reserve", Data: reserveColumn},
	}
	return
}

func calculateVaccineReserve(vaccinations []sciensano.Vaccination, batches []*vaccines.Batch) (reserve []float64) {
	batchIndex := 0
	lastBatch := 0

	for _, entry := range vaccinations {
		// find the last time we received vaccines
		for batchIndex < len(batches) &&
			!batches[batchIndex].Date.Time.After(entry.Timestamp) {
			// how many vaccines have we received so far?
			lastBatch = batches[batchIndex].Amount
			batchIndex++
		}

		// add it to the list
		reserve = append(reserve, float64(lastBatch-entry.Total()))
	}

	return
}

func (handler *Handler) buildVaccineTimeTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var batches []*vaccines.Batch
	batches, err = handler.Vaccines.GetBatches(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccine data: %s", err.Error())
	}

	batches = vaccines.AccumulateBatches(batches)

	var vaccinations []sciensano.Vaccination
	vaccinations, err = handler.Sciensano.GetVaccinations(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("failed to get vaccination data: %s", err.Error())
	}

	vaccinations = sciensano.AccumulateVaccinations(vaccinations)
	timestampColumn, timeColumn := CalculateVaccineDelay(vaccinations, batches)

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestampColumn},
		{Text: "time", Data: timeColumn},
	}
	return
}

func CalculateVaccineDelay(vaccinations []sciensano.Vaccination, batches []*vaccines.Batch) (timestamps grafanajson.TableQueryResponseTimeColumn, delays grafanajson.TableQueryResponseNumberColumn) {
	batchIndex := 0
	batchCount := len(batches)

	for _, entry := range vaccinations {
		// how many vaccines did we need to perform this many vaccinations?
		vaccinesNeeded := entry.Partial + entry.Full

		// find when we reached that number of vaccines
		for batchIndex < batchCount && batches[batchIndex].Amount < vaccinesNeeded {
			batchIndex++
		}

		// we depleted the *previous* batch. report the time difference between now and when we received that batch
		if batchIndex > 0 {
			timestamps = append(timestamps, entry.Timestamp)
			delays = append(delays, entry.Timestamp.Sub(batches[batchIndex-1].Date.Time).Hours()/24)
		}
	}
	return
}
