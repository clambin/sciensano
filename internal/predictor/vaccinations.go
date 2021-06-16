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

	firstDoses := make(chan []float64)
	secondDoses := make(chan []float64)
	output := make(chan []float64)

	input := buildVaccinationInput(vaccinations /*[len(vaccinations)-HistoryBatches*BatchSize:]*/)

	go forecastSamples(input[0], input[1], ForecastBatches*BatchSize, "first vaccination", firstDoses)
	go forecastSamples(input[1], input[0], ForecastBatches*BatchSize, "second vaccination", secondDoses)

	score1 := <-firstDoses
	score2 := <-secondDoses

	score = score1[0] * score2[0]

	const (
		a = 10
		b = 0
	)
	go consolidateSamples(output, []chan []float64{firstDoses, secondDoses}, func(input [][]float64) []float64 {
		firstDose := math.Max(0.0, (a*input[0][0]+b*input[1][1])/(a+b))
		secondDose := math.Max(0.0, (b*input[0][1]+a*input[1][0])/(a+b))

		return []float64{firstDose, secondDose}
	})

	begin, _, delta := getVaccinationDates(vaccinations)
	end := begin.Add(BatchSize * delta)

	for values := range output {
		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  end,
			FirstDose:  int(values[0]),
			SecondDose: int(math.Min(values[0], values[1])),
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
