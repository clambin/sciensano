package apihandler

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"time"
)

type VaccinationLag struct {
	Timestamp time.Time
	FullDose  int
	Lag       float64

	index int
}

func buildLag(vaccinations []sciensano.Vaccination) (lag []VaccinationLag) {
	// record all full vaccinations
	var (
		currentFull, firstAtFull, index int
		vaccination                     sciensano.Vaccination
		full                            VaccinationLag
	)

	for index, vaccination = range vaccinations {
		if vaccination.SecondDose != 0 && (vaccination.SecondDose > currentFull || vaccination.FirstDose != firstAtFull) {
			lag = append(lag, VaccinationLag{
				Timestamp: vaccination.Timestamp,
				FullDose:  vaccination.SecondDose,
				index:     index,
			})
			currentFull = vaccination.SecondDose
			firstAtFull = vaccination.FirstDose
		}
	}

	for index, full = range lag {
		// find the time when firstDose equals secondDose
		// we may not find any occurrence of when firstDose was the recorded lastDose (initial data may be complete).
		// don't report a delta larger than vs. the first recorded vaccination
		lastTime := lag[0].Timestamp
		var index2 int
		for index2, vaccination = range vaccinations {
			if vaccination.FirstDose <= full.FullDose && index2 <= full.index {
				lastTime = vaccination.Timestamp
			} else {
				break
			}
		}
		lag[index].Lag = full.Timestamp.Sub(lastTime).Hours() / 24.0
	}
	return
}
