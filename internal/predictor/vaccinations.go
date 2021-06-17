package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	"math"
	"time"
)

func ForecastVaccinations(vaccinations []sciensano.Vaccination) (forecast []sciensano.Vaccination, err error) {
	if len(vaccinations) < BatchSize {
		return nil, fmt.Errorf("not enough data: at least %d samples required", BatchSize)
	}

	input := buildVaccinationInput(vaccinations)

	firstDoses := forecastSamples([][]float64{input[0]}, ForecastSamples, "first vaccination")
	secondDoses := forecastSamples([][]float64{input[1]}, ForecastSamples, "second vaccination")
	output := consolidateSamples([]chan []float64{firstDoses, secondDoses}, singleConsolidator)

	begin, _, delta := getVaccinationDates(vaccinations)
	end := begin.Add(BatchSize * delta)

	for values := range output {
		first := values[0]
		ratio := math.Min(1.0, values[1])
		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  end,
			FirstDose:  int(first),
			SecondDose: int(math.Min(first, ratio*first)),
		})

		end = end.Add(delta)
	}

	return
}

func buildVaccinationInput(vaccinations []sciensano.Vaccination) (output [][]float64) {
	output = make([][]float64, 2)
	for _, test := range vaccinations {
		output[0] = append(output[0], float64(test.FirstDose))
		if test.FirstDose > 0 {
			output[1] = append(output[1], float64(test.SecondDose)/float64(test.FirstDose))
		} else {
			output[1] = append(output[1], 0.0)
		}
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
