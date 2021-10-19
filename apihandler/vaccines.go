package apihandler

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	"sync"
	"time"
)

func (handler *Handler) buildVaccineTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse) {
	if batches, err := handler.Vaccines.GetBatches(ctx); err == nil {
		batches = vaccines.AccumulateBatches(batches)

		rows := len(batches)
		timestampColumn := make(grafanaJson.TableQueryResponseTimeColumn, 0, rows)
		batchColumn := make(grafanaJson.TableQueryResponseNumberColumn, 0, rows)

		for _, entry := range batches {
			if entry.Date.Time.After(endTime) {
				continue
			}
			timestampColumn = append(timestampColumn, entry.Date.Time)
			batchColumn = append(batchColumn, float64(entry.Amount))
		}

		response = new(grafanaJson.TableQueryResponse)
		response.Columns = []grafanaJson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "vaccines", Data: batchColumn},
		}
	}
	return
}

func (handler *Handler) buildVaccineStatsTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse) {
	var batches []*vaccines.Batch
	var vaccinations []sciensano.Vaccination
	var err error
	if batches, err = handler.Vaccines.GetBatches(ctx); err == nil {
		batches = vaccines.AccumulateBatches(batches)
		if vaccinations, err = handler.Sciensano.GetVaccinations(ctx, endTime); err == nil {
			vaccinations = sciensano.AccumulateVaccinations(vaccinations)
		}

		rows := len(vaccinations)
		timestampColumn := make(grafanaJson.TableQueryResponseTimeColumn, 0, rows)
		vaccinationsColumn := make(grafanaJson.TableQueryResponseNumberColumn, 0, rows)
		reserveColumn := make(grafanaJson.TableQueryResponseNumberColumn, 0, rows)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			for _, entry := range vaccinations {
				timestampColumn = append(timestampColumn, entry.Timestamp)
				vaccinationsColumn = append(vaccinationsColumn, float64(entry.Partial+entry.Full))
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

		response = new(grafanaJson.TableQueryResponse)
		response.Columns = []grafanaJson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "vaccinations", Data: vaccinationsColumn},
			{Text: "reserve", Data: reserveColumn},
		}
	}
	return
}

func (handler *Handler) buildVaccineTimeTableResponse(ctx context.Context, _, endTime time.Time, _ string) (response *grafanaJson.TableQueryResponse) {
	var batches []*vaccines.Batch
	var vaccinations []sciensano.Vaccination
	var err error
	if batches, err = handler.Vaccines.GetBatches(ctx); err == nil {
		batches = vaccines.AccumulateBatches(batches)
		if vaccinations, err = handler.Sciensano.GetVaccinations(ctx, endTime); err == nil {
			vaccinations = sciensano.AccumulateVaccinations(vaccinations)
		}

		timestampColumn, timeColumn := CalculateVaccineDelay(vaccinations, batches)

		response = new(grafanaJson.TableQueryResponse)
		response.Columns = []grafanaJson.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "time", Data: timeColumn},
		}
	}
	return
}

func CalculateVaccineDelay(vaccinations []sciensano.Vaccination, batches []*vaccines.Batch) (timestamps grafanaJson.TableQueryResponseTimeColumn, delays grafanaJson.TableQueryResponseNumberColumn) {
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
		reserve = append(reserve, float64(lastBatch-entry.Full-entry.Partial))
	}

	return
}
