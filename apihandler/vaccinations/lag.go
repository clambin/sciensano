package vaccinations

import (
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/datasets"
)

func buildLag(vaccinationsData *datasets.Dataset) (timestamps grafanaJson.TableQueryResponseTimeColumn, lag grafanaJson.TableQueryResponseNumberColumn) {
	var firstDoseIndex, lastSecondDose int

	/*
		if len(vaccinationsData.Groups) == 0 {
			log.Warning("no vaccination data to calculate lag")
			return
		}
	*/

	// run through all vaccinations
	for index, entry := range vaccinationsData.Groups[0].Values {
		// we only measure lag when there is actually a second dose
		// we don't report when the 2nd dose doesn't change
		if entry.(*sciensano.VaccinationsEntry).Full == 0 || entry.(*sciensano.VaccinationsEntry).Full == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals (or higher) the current Second Dose number
		for firstDoseIndex <= index && vaccinationsData.Groups[0].Values[firstDoseIndex].(*sciensano.VaccinationsEntry).Partial < entry.(*sciensano.VaccinationsEntry).Full {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, vaccinationsData.Timestamps[index])
			lag = append(lag, vaccinationsData.Timestamps[index].Sub(vaccinationsData.Timestamps[firstDoseIndex]).Hours()/24)
		}

		lastSecondDose = entry.(*sciensano.VaccinationsEntry).Full
	}

	return
}
