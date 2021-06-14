package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	"math"
	"time"
)

func ForecastVaccinations(vaccinations []sciensano.Vaccination) (forecast []sciensano.Vaccination, score float64, err error) {
	if len(vaccinations) < BatchSize {
		return nil, 0.0, fmt.Errorf("not enough data: at least %d samples required", BatchSize)
	}

	firstDoses := make(chan float64)
	secondDoses := make(chan float64)

	input := buildVaccinationInput(vaccinations)

	go forecastSamples(input[0], input[1], ForecastBatches*BatchSize, "vaccination", firstDoses)
	go forecastSamples(input[1], input[0], ForecastBatches*BatchSize, "vaccination", secondDoses)

	_, end, delta := getVaccinationDates(vaccinations)

	score1 := <-firstDoses
	score2 := <-secondDoses

	score = score1 * score2

	for i := 0; i < ForecastBatches*BatchSize; i++ {
		end = end.Add(delta)

		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  end,
			FirstDose:  int(<-firstDoses),
			SecondDose: int(<-secondDoses),
		})
	}

	return
}

func buildVaccinationInput(vaccinations []sciensano.Vaccination) (output [][]float64) {
	output = make([][]float64, 2)
	for _, test := range vaccinations[len(vaccinations)-HistoryBatches*BatchSize:] {
		output[0] = append(output[0], float64(test.FirstDose))
		output[1] = append(output[1], float64(test.SecondDose))
	}
	return
}

func getVaccinationDates(vaccinations []sciensano.Vaccination) (from, to time.Time, delta time.Duration) {
	from = vaccinations[0].Timestamp
	to = vaccinations[len(vaccinations)-1].Timestamp
	avg := math.Round(to.Sub(from).Hours() / float64(len(vaccinations)))
	delta = time.Duration(avg) * time.Hour
	return
}
