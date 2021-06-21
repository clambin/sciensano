package predictor

import (
	log "github.com/sirupsen/logrus"
)

const (
	BatchSize           = 7
	ForecastSampleCount = 3 * BatchSize
	learnThreshold      = 0.945
	learnRetries        = 1
)

func ForecastSamples(forecastCount, batchSize int, label string, series ...[]float64) (output chan []float64) {
	output = make(chan []float64)
	go forecast(output, forecastCount, batchSize, label, series)
	return
}

func forecast(responses chan []float64, forecastCount, batchSize int, label string, series [][]float64) {
	p := New(batchSize, 1000)

	totalSeries := make([][]float64, 0)

	for _, entry := range series {
		if len(entry) != 0 {
			totalSeries = append(totalSeries, entry)
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

	buffer := make([]float64, batchSize)
	copy(buffer, series[0][:batchSize])

	for i := batchSize; i < len(totalSeries[0])+forecastCount; i++ {
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

func ConsolidateSamples(processor func([][]float64) []float64, input ...chan []float64) (output chan []float64) {
	output = make(chan []float64)
	go consolidate(output, processor, input)
	return
}

func consolidate(output chan []float64, processor func([][]float64) []float64, input []chan []float64) {
	allValues := make([][]float64, len(input))

	for values := range input[0] {
		allValues[0] = values

		for index, channel := range input[1:] {
			allValues[index+1] = <-channel
		}

		output <- processor(allValues)
	}

	close(output)
}

/*
func StandardConsolidator(input [][]float64) []float64 {
	const a = 1
	const b = 0

	return []float64{
		math.Max(0.0, (a*input[0][0]+b*input[1][1])/(a+b)),
		math.Max(0.0, (b*input[0][1]+a*input[1][0])/(a+b)),
	}
}
*/

func SingleConsolidator(input [][]float64) []float64 {
	return []float64{input[0][0], input[1][0]}
}
