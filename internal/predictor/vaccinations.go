package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

func ForecastVaccinations(vaccinations []sciensano.Vaccination) (forecast []sciensano.Vaccination, score float64, err error) {
	partialsResponse := make(chan struct {
		forecast []float64
		score    float64
		err      error
	})
	go func() {
		resp := struct {
			forecast []float64
			score    float64
			err      error
		}{}
		resp.forecast, resp.score, resp.err = forecastVaccinations(vaccinations, func(vaccination sciensano.Vaccination) int { return vaccination.FirstDose })
		partialsResponse <- resp
	}()

	fullResponse := make(chan struct {
		forecast []float64
		score    float64
		err      error
	})
	go func() {
		resp := struct {
			forecast []float64
			score    float64
			err      error
		}{}
		resp.forecast, resp.score, resp.err = forecastVaccinations(vaccinations, func(vaccination sciensano.Vaccination) int { return vaccination.SecondDose })
		fullResponse <- resp
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

	for i := 0; i < len(partials.forecast); i++ {
		forecast = append(forecast, sciensano.Vaccination{
			Timestamp:  endDate,
			FirstDose:  int(partials.forecast[i]),
			SecondDose: int(full.forecast[i]),
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

	var i int
	for i = 0; score < 0.98 && i < 20; i++ {
		score = p.Learn(input)
	}

	log.WithField("score", score).Infof("learned from %d samples after %d attempts", len(input), i+1)

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
