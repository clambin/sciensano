package apihandler

import (
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/pkg/sciensano"
	"time"
)

type VaccinationLag struct {
	Timestamp time.Time
	FullDose  int
	Lag       float64

	index int
}

func buildLag(vaccinations []sciensano.Vaccination) (timestamps grafana_json.TableQueryResponseTimeColumn, lag grafana_json.TableQueryResponseNumberColumn) {
	// record all full vaccinations
	var (
		firstDoseIndex int
		lastSecondDose int
	)

	vaccinationCount := len(vaccinations)

	// run through all vaccinations
	for index := 0; index < vaccinationCount; index++ {
		// we only measure lag when there is actually a second dose
		if vaccinations[index].SecondDose == 0 {
			continue
		}
		// we don't report when the 2nd dose doesn't change
		if vaccinations[index].SecondDose == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals the current Second Dose number
		for firstDoseIndex <= index && vaccinations[firstDoseIndex].FirstDose < vaccinations[index].SecondDose {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, vaccinations[index].Timestamp)
			lag = append(lag, vaccinations[index].Timestamp.Sub(vaccinations[firstDoseIndex].Timestamp).Hours()/24)
		}

		lastSecondDose = vaccinations[index].SecondDose
	}

	return
}
