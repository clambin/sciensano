package apihandler

import (
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"time"
)

type VaccinationLag struct {
	Timestamp time.Time
	FullDose  int
	Lag       float64

	index int
}

func buildLag(vaccinations []sciensano.Vaccination) (timestamps grafanaJson.TableQueryResponseTimeColumn, lag grafanaJson.TableQueryResponseNumberColumn) {
	// record all full vaccinations
	var (
		firstDoseIndex int
		lastSecondDose int
	)

	// run through all vaccinations
	for index, entry := range vaccinations {
		// we only measure lag when there is actually a second dose
		// we don't report when the 2nd dose doesn't change
		if entry.Full == 0 || entry.Full == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals (or higher) the current Second Dose number
		for firstDoseIndex <= index && vaccinations[firstDoseIndex].Partial < entry.Full {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, entry.Timestamp)
			lag = append(lag, entry.Timestamp.Sub(vaccinations[firstDoseIndex].Timestamp).Hours()/24)
		}

		lastSecondDose = vaccinations[index].Full
	}

	return
}
