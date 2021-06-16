package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	"math"
	"time"
)

func ForecastTests(tests []sciensano.Test) (forecast []sciensano.Test, score float64, err error) {
	if len(tests) < HistoryBatches*BatchSize {
		return nil, 0.0, fmt.Errorf("not enough data: at least %d samples required", HistoryBatches*BatchSize)
	}

	totalTests := make(chan float64)
	positiveTests := make(chan float64)

	input := buildTestsInput(tests[len(tests)-HistoryBatches*BatchSize:])

	// sklearn doesn't give us forecasts for both data sets in one prediction (gonum doesn't support n-dimensional arrays),
	// so we run both forecasts in parallel.  we're still passing both streams to train the models for both input streams
	// (though not really clear if this works at all ...)

	go forecastSamples(input[0], input[1], ForecastBatches*BatchSize, "total test", totalTests)
	go forecastSamples(input[1], input[0], ForecastBatches*BatchSize, "positive test", positiveTests)

	_, end, delta := getTestDates(tests)

	score1 := <-totalTests
	score2 := <-positiveTests

	score = score1 * score2

	for i := 0; i < ForecastBatches*BatchSize; i++ {
		end = end.Add(delta)

		total1 := <-totalTests
		positive1 := <-totalTests
		positive2 := <-positiveTests
		total2 := <-positiveTests

		total := (8*total1 + 2*total2) / 10
		positive := (2*positive1 + 8*positive2) / 10

		forecast = append(forecast, sciensano.Test{
			Timestamp: end,
			Total:     int(math.Max(0.0, total)),
			Positive:  int(math.Max(0.0, positive)),
		})
	}

	return
}

func buildTestsInput(tests []sciensano.Test) (output [][]float64) {
	output = make([][]float64, 2)
	for _, test := range tests {
		output[0] = append(output[0], float64(test.Total))
		output[1] = append(output[1], float64(test.Positive))
	}
	return
}

func getTestDates(tests []sciensano.Test) (from, to time.Time, delta time.Duration) {
	from = tests[0].Timestamp
	to = tests[len(tests)-1].Timestamp
	avg := math.Round(to.Sub(from).Hours() / float64(len(tests)))
	delta = time.Duration(avg) * time.Hour
	return
}
