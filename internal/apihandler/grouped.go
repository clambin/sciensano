package apihandler

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"sort"
	"time"
)

func getTimestamps(vaccinations map[string][]sciensano.Vaccination) (timestamps []time.Time) {
	// get unique timestamps
	uniqueTimestamps := make(map[time.Time]bool, len(vaccinations))
	for _, groupData := range vaccinations {
		for _, data := range groupData {
			uniqueTimestamps[data.Timestamp] = true
		}
	}
	timestamps = make([]time.Time, 0, len(uniqueTimestamps))
	for timestamp := range uniqueTimestamps {
		timestamps = append(timestamps, timestamp)
	}
	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i].Before(timestamps[j]) })
	return
}

// getFilledVaccinations has two main goals: 1) it returns either the partial or complete vaccination figures for a group and
// 2) it fills out the series with any missing timestamps so all columns cover the complete range of timestamps
func getFilledVaccinations(timestamps []time.Time, vaccinations []sciensano.Vaccination, complete bool) (filled []float64) {
	timestampCount := len(timestamps)
	vaccinationCount := len(vaccinations)

	var timestampIndex, vaccinationIndex int
	var lastVaccination sciensano.Vaccination

	for timestampIndex < timestampCount {
		for vaccinationIndex < vaccinationCount && timestamps[timestampIndex].Before(vaccinations[vaccinationIndex].Timestamp) {
			lastVaccination.Timestamp = timestamps[timestampIndex]
			filled = append(filled, float64(getVaccination(lastVaccination, complete)))
			timestampIndex++
		}
		if vaccinationIndex < vaccinationCount && timestamps[timestampIndex].Equal(vaccinations[vaccinationIndex].Timestamp) {
			lastVaccination = vaccinations[vaccinationIndex]
			filled = append(filled, float64(getVaccination(lastVaccination, complete)))
			vaccinationIndex++
			timestampIndex++
		} else if vaccinationIndex == vaccinationIndex {
			for ; timestampIndex < timestampCount; timestampIndex++ {
				lastVaccination.Timestamp = timestamps[timestampIndex]
				filled = append(filled, float64(getVaccination(lastVaccination, complete)))
			}
		}
	}
	return
}

func getVaccination(vaccination sciensano.Vaccination, complete bool) (value int) {
	if complete == false {
		value = vaccination.FirstDose
	} else {
		value = vaccination.SecondDose
	}
	return
}
