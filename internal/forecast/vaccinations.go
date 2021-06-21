package forecast

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"math"
	"time"
)

func PredictVaccinations(vaccinations []sciensano.Vaccination) (forecast []sciensano.Vaccination, err error) {
	if len(vaccinations) < predictor.BatchSize {
		return nil, fmt.Errorf("not enough data: at least %d samples required", predictor.BatchSize)
	}

	input := buildVaccinationInput(vaccinations)

	firstDoses := predictor.ForecastSamples(predictor.ForecastSampleCount, predictor.BatchSize, "first vaccination", input[0])
	secondDoses := predictor.ForecastSamples(predictor.ForecastSampleCount, predictor.BatchSize, "second vaccination", input[1])
	output := predictor.ConsolidateSamples(predictor.SingleConsolidator, firstDoses, secondDoses)

	begin, _, delta := getVaccinationDates(vaccinations)
	end := begin.Add(predictor.BatchSize * delta)

	for values := range output {
		first := values[0]
		second := math.Min(first, values[1])
		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  end,
			FirstDose:  int(first),
			SecondDose: int(second),
		})

		end = end.Add(delta)
	}

	return
}

func buildVaccinationInput(vaccinations []sciensano.Vaccination) (output [][]float64) {
	output = make([][]float64, 2)
	for _, test := range vaccinations {
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
