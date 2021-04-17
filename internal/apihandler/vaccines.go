package apihandler

import (
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/pkg/sciensano"
	"time"
)

func calculateVaccineDelay(vaccinations []sciensano.Vaccination, batches []vaccines.Batch) (timestamps []time.Time, delays []float64) {
	var batchIndex int

	for _, entry := range vaccinations {
		// how many vaccines did we need to perform this many vaccinations?
		vaccinesNeeded := entry.FirstDose + entry.SecondDose

		// find when we reached that number of vaccines
		for batchIndex < len(batches) &&
			batches[batchIndex].Amount < vaccinesNeeded {
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
