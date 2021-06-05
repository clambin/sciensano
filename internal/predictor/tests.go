package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

const (
	batchSize       = 56
	historyBatches  = 8
	forecastBatches = 1
)

func getDates(tests []sciensano.Test) (from, to time.Time, delta time.Duration) {
	from = tests[0].Timestamp
	to = tests[len(tests)-1].Timestamp
	avg := math.Round(to.Sub(from).Hours() / float64(len(tests)))
	delta = time.Duration(avg) * time.Hour
	return
}

func ForecastTests(tests []sciensano.Test) (forecast []sciensano.Test, err error) {
	if len(tests) < batchSize*historyBatches {
		return nil, fmt.Errorf("not enough data: at least %d samples required", batchSize*historyBatches)
	}

	var totals []float64
	totals, err = forecastTestsAttribute(tests, func(test sciensano.Test) float64 { return float64(test.Total) })

	if err != nil {
		return nil, err
	}

	var positives []float64
	positives, err = forecastTestsAttribute(tests, func(test sciensano.Test) float64 { return float64(test.Positive) })

	_, end, delta := getDates(tests)

	for index := range totals {
		end = end.Add(delta)
		forecast = append(forecast, sciensano.Test{
			Timestamp: end,
			Total:     int(math.Max(0, totals[index])),
			Positive:  int(math.Max(0, positives[index])),
		})
	}

	return
}

func forecastTestsAttribute(tests []sciensano.Test, attribute func(test sciensano.Test) float64) (forecast []float64, err error) {
	p := New(batchSize, 100000)

	input := make([]float64, len(tests))
	for i, test := range tests {
		input[i] = attribute(test)
	}

	score := 0.0
	for i := 0; score < 0.80 && i < 20; i++ {
		score = p.Learn(input)
	}
	log.WithField("score", score).Infof("analyzing %d samples for total tests", len(input))

	output := make([]float64, batchSize)
	copy(output, input[len(input)-batchSize:])

	for i := 0; i < forecastBatches; i++ {
		output, err = p.PredictN(output, batchSize)

		if err == nil {
			forecast = append(forecast, output...)
		}
	}

	return
}
