package apihandler

import (
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/pkg/sciensano"
	"time"
)

func calculateVaccineDelay(vaccinations []sciensano.Vaccination, batches []vaccines.Batch) (delays []float64) {
	var batchIndex int

	for _, entry := range vaccinations {
		// how many vaccines did we need to perform this many vaccinations?
		vaccinesNeeded := entry.FirstDose + entry.SecondDose

		// find when we reached the number of vaccines to perform this number
		for batchIndex < len(batches) &&
			batches[batchIndex].Amount < vaccinesNeeded {
			batchIndex++
		}

		// add it to the list
		if batchIndex < len(batches) {
			delay := entry.Timestamp.Sub(time.Time(batches[batchIndex].Date)).Hours() / 24
			delays = append(delays, delay)
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
