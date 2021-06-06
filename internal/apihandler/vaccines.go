package apihandler

import (
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/pkg/sciensano"
	"sync"
	"time"
)

func (handler *Handler) buildVaccineTableResponse(_ time.Time, _ string) (response *grafana_json.TableQueryResponse) {
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

func (handler *Handler) buildVaccineStatsTableResponse(endTime time.Time, _ string) (response *grafana_json.TableQueryResponse) {
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
		vaccinationsColumn := make(grafana_json.TableQueryResponseNumberColumn, 0, rows)
		reserveColumn := make(grafana_json.TableQueryResponseNumberColumn, 0, rows)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			for _, entry := range vaccinations {
				timestampColumn = append(timestampColumn, entry.Timestamp)
				vaccinationsColumn = append(vaccinationsColumn, float64(entry.FirstDose+entry.SecondDose))
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

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "vaccinations", Data: vaccinationsColumn},
			{Text: "reserve", Data: reserveColumn},
		}
	}
	return
}

func (handler *Handler) buildVaccineTimeTableResponse(endTime time.Time, _ string) (response *grafana_json.TableQueryResponse) {
	var batches []vaccines.Batch
	var vaccinations []sciensano.Vaccination
	var err error
	if batches, err = handler.Vaccines.GetBatches(); err == nil {
		batches = vaccines.AccumulateBatches(batches)
		if vaccinations, err = handler.Sciensano.GetVaccinations(endTime); err == nil {
			vaccinations = sciensano.AccumulateVaccinations(vaccinations)
		}

		timestampColumn, timeColumn := CalculateVaccineDelay(vaccinations, batches)

		response = new(grafana_json.TableQueryResponse)
		response.Columns = []grafana_json.TableQueryResponseColumn{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "time", Data: timeColumn},
		}
	}
	return
}

func CalculateVaccineDelay(vaccinations []sciensano.Vaccination, batches []vaccines.Batch) (timestamps grafana_json.TableQueryResponseTimeColumn, delays grafana_json.TableQueryResponseNumberColumn) {
	batchIndex := 0
	batchCount := len(batches)

	for _, entry := range vaccinations {
		// how many vaccines did we need to perform this many vaccinations?
		vaccinesNeeded := entry.FirstDose + entry.SecondDose

		// find when we reached that number of vaccines
		for batchIndex < batchCount && batches[batchIndex].Amount < vaccinesNeeded {
			batchIndex++
		}

		// we depleted the *previous* batch. report the time difference between now and when we received that batch
		if batchIndex > 0 {
			timestamps = append(timestamps, entry.Timestamp)
			delays = append(delays, entry.Timestamp.Sub(time.Time(batches[batchIndex-1].Date)).Hours()/24)
		}
	}
	return
}

func calculateVaccineReserve(vaccinations []sciensano.Vaccination, batches []vaccines.Batch) (reserve []float64) {
	batchIndex := 0
	lastBatch := 0

	for _, entry := range vaccinations {
		// find the last time we received vaccines
		for batchIndex < len(batches) &&
			!time.Time(batches[batchIndex].Date).After(entry.Timestamp) {
			// how many vaccines have we received so far?
			lastBatch = batches[batchIndex].Amount
			batchIndex++
		}

		// add it to the list
		reserve = append(reserve, float64(lastBatch-entry.SecondDose-entry.FirstDose))
	}

	return
}
