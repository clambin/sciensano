package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

func ForecastVaccinations(vaccinations []sciensano.Vaccination) (forecast []sciensano.Vaccination, score float64, err error) {
	var partials []float64
	var score1 float64

	partials, score1, err = forecastVaccinations(vaccinations, func(vaccination sciensano.Vaccination) int { return vaccination.FirstDose })

	if err != nil {
		return nil, 0, err
	}

	var full []float64
	var score2 float64
	full, score2, err = forecastVaccinations(vaccinations, func(vaccination sciensano.Vaccination) int { return vaccination.SecondDose })

	if err != nil {
		return nil, 0, err
	}

	score = (score1 + score2) / 2

	_, endDate, delta := getVaccinationDates(vaccinations)

	for i := 0; i < len(partials); i++ {
		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  endDate,
			FirstDose:  int(partials[i]),
			SecondDose: int(full[i]),
		})
		endDate = endDate.Add(delta)
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

func forecastVaccinations(vaccinations []sciensano.Vaccination, getAttribute func(vaccination sciensano.Vaccination) int) (forecast []float64, score float64, err error) {
	if len(vaccinations) < batchSize {
		return nil, 0.0, fmt.Errorf("need at least %d samples", batchSize)
	}

	input := make([]float64, len(vaccinations))
	for index, vaccination := range vaccinations {
		input[index] = float64(getAttribute(vaccination))
	}

	p := New(batchSize, 10000)

	for score < 0.94 {
		score = p.Learn(input)
	}

	log.WithField("score", score).Infof("learned from %d samples", len(input))

	output := make([]float64, batchSize)
	copy(output, input[len(input)-batchSize:])

	lastValue := input[len(input)-1]

	for i := 0; err == nil && i < forecastBatches; i++ {
		if output, err = p.PredictN(output, batchSize); err == nil {
			for _, value := range output {
				if value > lastValue {
					forecast = append(forecast, value)
					lastValue = value
				} else {
					forecast = append(forecast, lastValue)
				}
			}
		}
	}

	return
}
