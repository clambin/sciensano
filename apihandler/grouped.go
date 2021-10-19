package apihandler

import (
	"github.com/clambin/sciensano/sciensano"
	"sort"
	"time"
)

func getTimestamps(vaccinations map[string][]sciensano.Vaccination) (timestamps []time.Time) {
	// get unique timestamps
	uniqueTimestamps := make(map[time.Time]bool)
	for _, groupData := range vaccinations {
		for _, data := range groupData {
			if _, ok := uniqueTimestamps[data.Timestamp]; ok == false {
				uniqueTimestamps[data.Timestamp] = true
				timestamps = append(timestamps, data.Timestamp)
			}
		}
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

func fillVaccinations(timestamps []time.Time, vaccinations map[string][]sciensano.Vaccination, complete bool) (results map[string]chan []float64) {
	results = make(map[string]chan []float64)
	for group := range vaccinations {
		results[group] = make(chan []float64)
		go func(groupName string, channel chan []float64) {
			channel <- getFilledVaccinations(timestamps, vaccinations[groupName], complete)
		}(group, results[group])
	}
	return
}

func getVaccination(vaccination sciensano.Vaccination, complete bool) (value int) {
	if complete == false {
		value = vaccination.Partial
	} else {
		value = vaccination.Full
	}
	return
}

func getGroups(vaccinations map[string][]sciensano.Vaccination) (groups []string) {
	for group := range vaccinations {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	return
}
