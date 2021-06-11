package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

func getDates(tests []sciensano.Test) (from, to time.Time, delta time.Duration) {
	from = tests[0].Timestamp
	to = tests[len(tests)-1].Timestamp
	avg := math.Round(to.Sub(from).Hours() / float64(len(tests)))
	delta = time.Duration(avg) * time.Hour
	return
}

func ForecastTests(tests []sciensano.Test) (forecast []sciensano.Test, score float64, err error) {
	if len(tests) < BatchSize {
		return nil, 0.0, fmt.Errorf("not enough data: at least %d samples required", BatchSize)
	}

	totalsResponse := make(chan forecastFigures)
	go func() {
		totalsResponse <- forecastTestsAttribute(tests, func(test sciensano.Test) float64 { return float64(test.Total) })
	}()

	positivesResponse := make(chan forecastFigures)
	go func() {
		positivesResponse <- forecastTestsAttribute(tests, func(test sciensano.Test) float64 { return float64(test.Positive) })
	}()

	totals := <-totalsResponse
	positives := <-positivesResponse

	if totals.err != nil {
		return nil, 0.0, totals.err
	} else if positives.err != nil {
		return nil, 0.0, positives.err
	}

	score = (totals.score + positives.score) / 2.0

	_, end, delta := getDates(tests)

	for index := range totals.figures {
		end = end.Add(delta)
		forecast = append(forecast, sciensano.Test{
			Timestamp: end,
			Total:     int(math.Max(0, totals.figures[index])),
			Positive:  int(math.Max(0, positives.figures[index])),
		})
	}

	return
}

func forecastTestsAttribute(tests []sciensano.Test, attribute func(test sciensano.Test) float64) (forecast forecastFigures) {
	if len(tests) < BatchSize {
		forecast.err = fmt.Errorf("need at least %d samples", BatchSize)
		return
	}

	input := make([]float64, len(tests))
	for i, test := range tests {
		input[i] = attribute(test)
	}

	p := New(BatchSize, 100000)

	for i := 0; forecast.score < 0.99 && i < learnRetries; i++ {
		forecast.score = p.Learn(input)
		log.WithField("score", forecast.score).Debugf("learned from %d test samples after %d attempts", len(input), 1+i)
	}

	output := make([]float64, BatchSize)
	copy(output, input[len(input)-BatchSize:])

	for i := 0; forecast.err == nil && i < ForecastBatches; i++ {
		output, forecast.err = p.PredictN(output, BatchSize)

		if forecast.err == nil {
			forecast.figures = append(forecast.figures, output...)
		}
	}

	return
}
