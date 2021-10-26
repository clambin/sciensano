package vaccinations

import (
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
)

func buildLag(vaccinationsData *sciensano.Vaccinations) (timestamps grafanaJson.TableQueryResponseTimeColumn, lag grafanaJson.TableQueryResponseNumberColumn) {
	// record all full vaccinations
	var (
		firstDoseIndex int
		lastSecondDose int
	)

	// run through all vaccinations
	for index, entry := range vaccinationsData.Groups[0].Values {
		// we only measure lag when there is actually a second dose
		// we don't report when the 2nd dose doesn't change
		if entry.Full == 0 || entry.Full == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals (or higher) the current Second Dose number
		for firstDoseIndex <= index && vaccinationsData.Groups[0].Values[firstDoseIndex].Partial < entry.Full {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, vaccinationsData.Timestamps[index])
			lag = append(lag, vaccinationsData.Timestamps[index].Sub(vaccinationsData.Timestamps[firstDoseIndex]).Hours()/24)
		}

		lastSecondDose = entry.Full
	}

	return
}
