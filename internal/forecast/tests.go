package forecast

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/predictor"
	"github.com/clambin/sciensano/pkg/sciensano"
	"math"
	"time"
)

const HistoryBatches = 30

func PredictTests(tests []sciensano.Test) (forecast []sciensano.Test, err error) {
	if len(tests) < predictor.BatchSize {
		return nil, fmt.Errorf("not enough data: at least %d samples required", predictor.BatchSize)
	}

	var history []sciensano.Test
	if len(tests) > HistoryBatches*predictor.BatchSize {
		history = tests[len(tests)-HistoryBatches*predictor.BatchSize:]
	} else {
		history = tests
	}

	input := buildTestsInput(history)

	// sklearn doesn't give us forecasts for both data sets in one prediction (gonum doesn't support n-dimensional arrays),
	// so we run both forecasts in parallel.
	//
	// tried training the two models for both datasets, but that didn't yield any improvements

	totalTests := predictor.ForecastSamples(predictor.ForecastSampleCount, predictor.BatchSize, "total test", input[0])
	positiveTests := predictor.ForecastSamples(predictor.ForecastSampleCount, predictor.BatchSize, "positive test", input[1])
	output := predictor.ConsolidateSamples(predictor.SingleConsolidator, totalTests, positiveTests)

	begin, _, delta := getTestDates(history)
	end := begin.Add(predictor.BatchSize * delta)

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
