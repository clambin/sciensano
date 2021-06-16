package predictor

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/sciensano"
	"math"
	"time"
)

const HistoryBatches = 30

func ForecastTests(tests []sciensano.Test) (forecast []sciensano.Test, err error) {
	if len(tests) < BatchSize {
		return nil, fmt.Errorf("not enough data: at least %d samples required", BatchSize)
	}

	history := tests
	if len(tests) > HistoryBatches*BatchSize {
		history = tests[len(tests)-HistoryBatches*BatchSize:]
	}

	totalTests := make(chan []float64)
	positiveTests := make(chan []float64)
	output := make(chan []float64)

	input := buildTestsInput(history)

	// sklearn doesn't give us forecasts for both data sets in one prediction (gonum doesn't support n-dimensional arrays),
	// so we run both forecasts in parallel.  we're still passing both streams to train the models for both input streams
	// (though not really clear if this works at all ...)

	go forecastSamples(input, ForecastSamples, "total test", totalTests)
	go forecastSamples([][]float64{input[1], input[0]}, ForecastSamples, "positive test", positiveTests)
	go consolidateSamples(output, []chan []float64{totalTests, positiveTests}, standardConsolidator)

	begin, _, delta := getTestDates(history)
	end := begin.Add(BatchSize * delta)

	for figures := range output {
		forecast = append(forecast, sciensano.Test{
			Timestamp: end,
			Total:     int(figures[0]),
			Positive:  int(math.Min(figures[1], figures[0])),
		})

		end = end.Add(delta)
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
