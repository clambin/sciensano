package predictor

import (
	log "github.com/sirupsen/logrus"
)

const (
	BatchSize       = 7
	ForecastSamples = 3 * BatchSize
	learnThreshold  = 0.945
	learnRetries    = 1
)

func forecastSamples(series [][]float64, forecastCount int, label string) (output chan []float64) {
	output = make(chan []float64)
	go forecast(output, series, forecastCount, label)
	return
}

func forecast(responses chan []float64, series [][]float64, forecastCount int, label string) {
	p := New(BatchSize, 1000)

	totalSeries := make([][]float64, 0)

	for _, serie := range series {
		if len(serie) != 0 {
			totalSeries = append(totalSeries, serie)
		}
	}

	var score float64
	retries := learnRetries
	for score < learnThreshold && retries > 0 {
		score = p.Learn(totalSeries)
		retries--
	}

	log.WithFields(log.Fields{
		"score":    score,
		"attempts": learnRetries - retries,
	}).Debugf("learned from %d %s samples", len(totalSeries[0]), label)

	buffer := make([]float64, BatchSize)
	copy(buffer, series[0][:BatchSize])

	for i := BatchSize; i < len(totalSeries[0])+forecastCount; i++ {
		prediction, _ := p.Predict(buffer)
		responses <- prediction

		if i < len(totalSeries[0]) {
			buffer = append(buffer[1:], totalSeries[0][i])
		} else {
			buffer = append(buffer[1:], prediction[0])
		}
	}

	close(responses)
}

func consolidateSamples(input []chan []float64, processor func([][]float64) []float64) (output chan []float64) {
	output = make(chan []float64)
	go consolidate(output, input, processor)
	return
}

func consolidate(output chan []float64, input []chan []float64, processor func([][]float64) []float64) {
	for values := range input[0] {
		allValues := [][]float64{values}

		for i := 1; i < len(input); i++ {
			allValues = append(allValues, <-input[i])
		}

		output <- processor(allValues)
	}

	close(output)
}

/*
func standardConsolidator(input [][]float64) []float64 {
	const a = 1
	const b = 0

	return []float64{
		math.Max(0.0, (a*input[0][0]+b*input[1][1])/(a+b)),
		math.Max(0.0, (b*input[0][1]+a*input[1][0])/(a+b)),
	}
}
*/

func singleConsolidator(input [][]float64) []float64 {
	return []float64{input[0][0], input[1][0]}
}
