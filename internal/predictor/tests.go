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

	totalTests := make(chan []float64)
	positiveTests := make(chan []float64)
	output := make(chan []float64)

	input := buildTestsInput(tests) // [len(tests)-HistoryBatches*BatchSize:])

	// sklearn doesn't give us forecasts for both data sets in one prediction (gonum doesn't support n-dimensional arrays),
	// so we run both forecasts in parallel.  we're still passing both streams to train the models for both input streams
	// (though not really clear if this works at all ...)

	go forecastSamples(input[0], input[1], ForecastBatches*BatchSize, "total test", totalTests)
	go forecastSamples(input[1], input[0], ForecastBatches*BatchSize, "positive test", positiveTests)

	score1 := <-totalTests
	score2 := <-positiveTests
	score = score1[0] * score2[0]

	const (
		a = 10
		b = 0
	)
	go consolidateSamples(output, []chan []float64{totalTests, positiveTests}, func(input [][]float64) []float64 {
		total := math.Max(0.0, (a*input[0][0]+b*input[1][1])/(a+b))
		positive := math.Max(0.0, (b*input[0][1]+a*input[1][0])/(a+b))

		return []float64{total, positive}
	})

	// _, end, delta := getTestDates(tests)
	begin, _, delta := getTestDates(tests)
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
