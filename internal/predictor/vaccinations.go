package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

func ForecastVaccinations(vaccinations []sciensano.Vaccination) (forecast []sciensano.Vaccination, score float64, err error) {
	partialsResponse := make(chan forecastFigures)
	fullResponse := make(chan forecastFigures)

	go func() {
		partialsResponse <- forecastVaccinations(vaccinations, func(vaccination sciensano.Vaccination) int { return vaccination.FirstDose })
	}()

	go func() {
		fullResponse <- forecastVaccinations(vaccinations, func(vaccination sciensano.Vaccination) int { return vaccination.SecondDose })
	}()

	partials := <-partialsResponse
	full := <-fullResponse

	if partials.err != nil {
		return nil, 0, partials.err
	}

	if full.err != nil {
		return nil, 0, full.err
	}

	score = (partials.score + full.score) / 2

	_, endDate, delta := getVaccinationDates(vaccinations)

	for i := 0; i < len(partials.figures); i++ {
		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  endDate,
			FirstDose:  int(partials.figures[i]),
			SecondDose: int(full.figures[i]),
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

func forecastVaccinations(vaccinations []sciensano.Vaccination, getAttribute func(vaccination sciensano.Vaccination) int) (forecast forecastFigures) {
	if len(vaccinations) < batchSize {
		forecast.err = fmt.Errorf("need at least %d samples", batchSize)
		return
	}

	input := make([]float64, len(vaccinations))
	for index, vaccination := range vaccinations {
		input[index] = float64(getAttribute(vaccination))
	}

	p := New(batchSize, 10000)

	var i int
	for i = 0; forecast.score < 0.99 && i < 20; i++ {
		forecast.score = p.Learn(input)
		log.WithField("score", forecast.score).Debugf("learned from %d vaccination samples after %d attempts", len(input), i+1)
	}

	output := make([]float64, batchSize)
	copy(output, input[len(input)-batchSize:])

	lastValue := input[len(input)-1]

	for i := 0; forecast.err == nil && i < forecastBatches; i++ {
		if output, forecast.err = p.PredictN(output, batchSize); forecast.err == nil {
			for _, value := range output {
				if value <= lastValue {
					value = lastValue
				}
				forecast.figures = append(forecast.figures, value)
				lastValue = value
			}
		}
	}

	return
}
