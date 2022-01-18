package vaccinations

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson"
)

func buildLag(vaccinationsData *datasets.Dataset) (timestamps simplejson.TableQueryResponseTimeColumn, lag simplejson.TableQueryResponseNumberColumn) {
	var firstDoseIndex, lastSecondDose int

	/*
		if len(vaccinationsData.Groups) == 0 {
			log.Warning("no vaccination data to calculate lag")
			return
		}
	*/

	// run through all vaccinations
	for index := range vaccinationsData.Timestamps {
		// we only measure lag when there is actually a second dose
		// we don't report when the 2nd dose doesn't change
		full := int(vaccinationsData.Groups[1].Values[index])
		if full == 0 || full == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals (or higher) the current Second Dose number
		for firstDoseIndex <= index && int(vaccinationsData.Groups[0].Values[firstDoseIndex]) < full {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, vaccinationsData.Timestamps[index])
			lag = append(lag, vaccinationsData.Timestamps[index].Sub(vaccinationsData.Timestamps[firstDoseIndex]).Hours()/24)
		}

		lastSecondDose = full
	}

	return
}
